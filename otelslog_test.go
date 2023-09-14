package otelslog_test

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/go-slog/otelslog"
)

func testHandler() slog.Handler {
	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     nil,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.TimeValue(time.Time{})
			}

			return a
		},
	})
}

func testContext() context.Context {
	ctx := context.Background()

	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: [16]byte{116, 114, 97, 99, 101, 95, 105, 100, 95, 116, 101, 115, 116, 49, 50, 51},
		SpanID:  [8]byte{115, 112, 97, 110, 95, 105, 100, 49},
	})
	ctx = trace.ContextWithSpanContext(ctx, spanContext)

	return ctx
}

func ExampleHandler() {
	var handler slog.Handler

	// Set up handler
	handler = testHandler()

	// Wrap handler
	handler = otelslog.NewHandler(handler)

	// Set up logger
	logger := slog.New(handler)

	// later in some trace context (eg. http request)

	ctx := testContext()

	// Call logger with a context
	logger.InfoContext(ctx, "hello world")

	// Output: time=0001-01-01T00:00:00.000Z level=INFO msg="hello world" trace_id=74726163655f69645f74657374313233 span_id=7370616e5f696431
}

func TestHandler_Handle_NoSpan(_ *testing.T) {
	handler := otelslog.NewHandler(slog.NewJSONHandler(io.Discard, nil))
	logger := slog.New(handler)

	logger.InfoContext(context.Background(), "hello world")
}
