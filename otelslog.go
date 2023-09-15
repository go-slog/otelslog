// Package otelslog provides an [slog.Handler] that attaches OpenTelemetry trace details to logs.
package otelslog

import (
	"context"
	"errors"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

const (
	traceIDKey = "trace_id"
	spanIDKey  = "span_id"
)

// NewHandler returns a new [Handler].
func NewHandler(handler slog.Handler) slog.Handler {
	return Handler{
		Handler: handler,
	}
}

// Middleware returns a [Middleware] for an [slogmulti.Pipe] handler.
//
// [Middleware]: https://pkg.go.dev/github.com/samber/slog-multi#Middleware
// [slogmulti.Pipe]: https://pkg.go.dev/github.com/samber/slog-multi#Pipe
func Middleware() func(slog.Handler) slog.Handler {
	return func(handler slog.Handler) slog.Handler {
		return NewHandler(handler)
	}
}

// Handler attaches details from an OpenTelemetry trace to each log record.
type Handler struct {
	slog.Handler
}

// Handle implements [slog.Handler].
func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	if h.Handler == nil {
		return errors.New("otelslog: handler is missing")
	}

	spanCtx := trace.SpanContextFromContext(ctx)

	if spanCtx.HasTraceID() {
		record.AddAttrs(slog.String(traceIDKey, spanCtx.TraceID().String()))
	}

	if spanCtx.HasSpanID() {
		record.AddAttrs(slog.String(spanIDKey, spanCtx.SpanID().String()))
	}

	return h.Handler.Handle(ctx, record)
}

// WithAttrs implements [slog.Handler].
func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if h.Handler == nil {
		return h
	}

	return Handler{h.Handler.WithAttrs(attrs)}
}

// WithGroup implements [slog.Handler].
func (h Handler) WithGroup(name string) slog.Handler {
	if h.Handler == nil {
		return h
	}

	return Handler{h.Handler.WithGroup(name)}
}
