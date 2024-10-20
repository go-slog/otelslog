package otelslog

import (
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

// ResourceMiddleware returns a [Middleware] for an [slogmulti.Pipe] handler attaching resource attributes to an [slog.Handler].
//
// [Middleware]: https://pkg.go.dev/github.com/samber/slog-multi#Middleware
// [slogmulti.Pipe]: https://pkg.go.dev/github.com/samber/slog-multi#Pipe
func ResourceMiddleware(res *resource.Resource) func(slog.Handler) slog.Handler {
	return func(handler slog.Handler) slog.Handler {
		return WithResource(handler, res)
	}
}

// WithResource returns an [slog.Handler] with all attributes attached from a [resource.Resource].
func WithResource(handler slog.Handler, res *resource.Resource) slog.Handler {
	if res.Len() == 0 {
		return handler
	}

	attrs := make([]slog.Attr, 0, res.Len())

	for iter := res.Iter(); iter.Next(); {
		curr := iter.Attribute()

		attrs = append(attrs, slog.Any(string(curr.Key), mapAttributeValue(curr.Value)))
	}

	return handler.WithAttrs(attrs)
}

func mapAttributeValue(v attribute.Value) slog.Value {
	switch v.Type() {
	case attribute.BOOL:
		return slog.BoolValue(v.AsBool())
	case attribute.INT64:
		return slog.Int64Value(v.AsInt64())
	case attribute.FLOAT64:
		return slog.Float64Value(v.AsFloat64())
	case attribute.STRING:
		return slog.StringValue(v.AsString())
	default:
		return slog.AnyValue(v.AsInterface())
	}
}
