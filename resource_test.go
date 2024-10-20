package otelslog_test

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/go-slog/otelslog"
)

func ExampleWithResource() {
	// Set up handler
	handler := testHandler()

	res, _ := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "my-service"),
			attribute.String("service.version", "1.0.0"),
			attribute.Bool("debug", false),
			attribute.Float64("sampling", 0.6),
		),
	)

	// Attach resource details to handler
	handler = otelslog.WithResource(handler, res)

	// Set up logger
	logger := slog.New(handler)

	// Call logger
	logger.Info("hello world")

	// Output: time=0001-01-01T00:00:00.000Z level=INFO msg="hello world" debug=false sampling=0.6 service.name=my-service service.version=1.0.0
}
