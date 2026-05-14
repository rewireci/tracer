package rewire

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// StartTestSpan creates a span named after t.Name() using the global tracer
// provider and registers t.Cleanup to end it. Call Init before using this if
// you want real export; without Init the global provider is a no-op and spans
// are silently discarded.
func StartTestSpan(t *testing.T) (context.Context, trace.Span) {
	t.Helper()
	ctx, span := otel.Tracer("rewire").Start(context.Background(), t.Name())
	t.Cleanup(func() { span.End() })
	return ctx, span
}
