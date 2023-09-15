package otelslog

import (
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

// WithResource returns an [slog.Logger] with all attributes attached from a [resource.Resource].
func WithResource(logger *slog.Logger, res *resource.Resource) *slog.Logger {
	if res.Len() == 0 {
		return logger
	}

	attrs := make([]any, 0, res.Len())

	for iter := res.Iter(); iter.Next(); {
		curr := iter.Attribute()

		attrs = append(attrs, slog.Any(string(curr.Key), mapAttributeValue(curr.Value)))
	}

	return logger.With(attrs...)
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
