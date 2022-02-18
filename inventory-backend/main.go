// inventory.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/albertteoh/gin-example/data"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
			semconv.ServiceNameKey.String("inventory-backend"),
		)),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	otel.SetTracerProvider(tp)
}

func getInventory() data.Inventory {
	return data.Inventory{
		Products: []data.Product{
			{Name: "potato", Price: 0.99, ID: "1"},
			{Name: "apple", Price: 0.50, ID: "2"},
			{Name: "mango", Price: 1.50, ID: "3"},
		},
	}
}

func main() {
	initTracer()

	r := gin.Default()
	r.Use(otelgin.Middleware("inventory-backend"))

	r.GET("/inventory", func(c *gin.Context) {
		prettyPrintHeader(c.Request.Header)
		c.JSON(200, gin.H{
			"inventory": getInventory(),
		})
	})
	r.Run("localhost:8081")
}

func prettyPrintHeader(header http.Header) {
	fmt.Printf("Request at %v\n", time.Now())
	for k, v := range header {
		fmt.Printf("%v: %v\n", k, v)
	}
}
