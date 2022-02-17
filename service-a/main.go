package main

import (
	"log"
	"net/http"

	"github.com/albertteoh/gin-example/data"
	"github.com/gin-gonic/gin"
	resty "github.com/go-resty/resty/v2"
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

func getInventory() data.Inventory {
	client := resty.New()
	otelTransport := HTTPClientTransporter(client.GetClient().Transport)
	client.SetTransport(otelTransport)
	var p Inventory
	// Call service-b
	_, err := client.R().
		SetResult(&p).
		Get("http://localhost:8081/inventory")

	if err != nil {
		log.Println("service-a: can't get inventory")
		return Inventory{}
	}
	log.Printf("service-a: go inventory")
	return p
}

func main() {
	initTracer()

	r := gin.Default()
	r.Use(otelgin.Middleware("inventory-server"))

	r.GET("/inventory", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"inventory": getInventory(),
		})
	})
	r.Run()
}
