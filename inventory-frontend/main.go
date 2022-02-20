package main

import (
	"context"
	"log"
	"net/http"

	"github.com/albertteoh/gin-example/data"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer() {
	otel.SetTextMapPropagator(propagation.TraceContext{})

	exp, err := jaeger.New(
		jaeger.WithAgentEndpoint(),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("inventory-frontend"),
		)),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	otel.SetTracerProvider(tp)
}

// HTTPClientTransporter is a convenience function which helps attaching tracing
// functionality to conventional HTTP clients.
func HTTPClientTransporter(rt http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(rt)
}

func getInventory(ctx context.Context) data.Inventory {
	client := resty.New()
	otelTransport := HTTPClientTransporter(client.GetClient().Transport)
	client.Debug = true
	client.SetTransport(otelTransport)
	var p data.Response
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetResult(&p).
		SetContext(ctx).
		Get("http://localhost:8082/inventory")

	if err != nil {
		log.Println("inventory-frontend: can't get inventory")
		return data.Inventory{}
	}
	log.Printf("inventory-frontend: got inventory %+v", p)
	return p.Inventory
}

func main() {
	initTracer()

	r := gin.Default()
	r.Use(otelgin.Middleware("inventory-frontend"))

	r.GET("/inventory", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"inventory": getInventory(c.Request.Context()),
		})
	})
	r.Run("localhost:8081")
}
