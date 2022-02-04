package tracing

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

type Tag = opentracing.Tag
type Tags = opentracing.Tags

type contextKey struct{}

var TracerKey = contextKey{}

/*
	Poor Mans OpenTracing.

	Standardizes logging of operation duration.
*/

type Tracer interface {
	StartSpan(operationName string) Span
}

type Span interface {
	SetBaggageItem(key string, value interface{})
	Finish()
}

// StartSpan starts and returns span with `operationName` and hooking as child to a span found within given context if any.
// It uses opentracing.Tracer propagated in context. If no found, it uses noop tracer without notification.
func StartSpan(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	tracer := tracerFromContext(ctx)
	if tracer == nil {
		// No tracing found, return noop span.
		return opentracing.NoopTracer{}.StartSpan(operationName), ctx
	}

	var span opentracing.Span
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	span = tracer.StartSpan(operationName, opts...)
	return span, opentracing.ContextWithSpan(ctx, span)
}

// ContextWithTracer returns a new `context.Context` that holds a reference to given opentracing.Tracer.
func ContextWithTracer(ctx context.Context, tracer opentracing.Tracer) context.Context {
	return context.WithValue(ctx, TracerKey, tracer)
}

// tracerFromContext extracts opentracing.Tracer from the given context.
func tracerFromContext(ctx context.Context) opentracing.Tracer {
	val := ctx.Value(TracerKey)
	if sp, ok := val.(opentracing.Tracer); ok {
		return sp
	}
	return nil
}
