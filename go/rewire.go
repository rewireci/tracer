// Package rewire provides OpenTelemetry instrumentation for CI pipelines.
// It auto-detects the CI platform and routes telemetry to Rewire (or any
// OTLP-compatible collector) using REWIRE_TOKEN or OTEL_EXPORTER_OTLP_ENDPOINT.
package rewire

import (
	"context"
	"log"
	"os"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const rewireEndpoint = "https://rewireci.com/otlp/v1"

var (
	mu         sync.Mutex
	shutdownFn func()
)

// Init configures the OpenTelemetry SDK from environment variables.
//
// Endpoint resolution order:
//  1. OTEL_EXPORTER_OTLP_ENDPOINT — used as-is (e.g. http://localhost:4318 from the Action)
//  2. REWIRE_TOKEN — sends directly to rewireci.com with Bearer auth
//  3. Neither set — tracing is disabled and a single warning is logged
//
// Returns a shutdown function that flushes pending spans and stops the tracer
// provider. Safe to call multiple times — subsequent calls return the same func.
func Init() func() {
	mu.Lock()
	defer mu.Unlock()
	if shutdownFn != nil {
		return shutdownFn
	}
	shutdownFn = initInternal()
	return shutdownFn
}

func initInternal() func() {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	token := os.Getenv("REWIRE_TOKEN")

	if endpoint == "" && token == "" {
		log.Println("[rewire] Neither OTEL_EXPORTER_OTLP_ENDPOINT nor REWIRE_TOKEN is set — tracing disabled")
		return func() {}
	}

	ci := DetectCI()

	attrs := []attribute.KeyValue{
		attribute.String("service.name", serviceName()),
		attribute.String("ci.platform", ci.Platform),
	}
	if ci.RunID != "" {
		attrs = append(attrs, attribute.String("run.id", ci.RunID))
	}

	res := resource.NewWithAttributes("", attrs...)

	// When OTEL_EXPORTER_OTLP_ENDPOINT is set the SDK reads it automatically;
	// no explicit option needed. When only REWIRE_TOKEN is set, configure the
	// Rewire endpoint and auth header explicitly.
	var opts []otlptracehttp.Option
	if endpoint == "" {
		opts = append(opts,
			otlptracehttp.WithEndpointURL(rewireEndpoint+"/traces"),
			otlptracehttp.WithHeaders(map[string]string{"Authorization": "Bearer " + token}),
		)
	}

	exporter, err := otlptracehttp.New(context.Background(), opts...)
	if err != nil {
		log.Printf("[rewire] Failed to create OTLP exporter: %v — tracing disabled\n", err)
		return func() {}
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("[rewire] Shutdown error: %v\n", err)
		}
	}
}

func serviceName() string {
	if v := os.Getenv("OTEL_SERVICE_NAME"); v != "" {
		return v
	}
	if v := os.Getenv("GITHUB_REPOSITORY"); v != "" {
		return v
	}
	return "unknown"
}

// Shutdown calls the shutdown function registered by Init.
// Safe to call before Init — does nothing in that case.
func Shutdown() {
	mu.Lock()
	fn := shutdownFn
	mu.Unlock()
	if fn != nil {
		fn()
	}
}

// _reset shuts down any active tracer provider and clears cached state.
// For use in tests only.
func _reset() {
	mu.Lock()
	fn := shutdownFn
	shutdownFn = nil
	mu.Unlock()

	if fn != nil {
		fn()
	}
	otel.SetTracerProvider(trace.NewNoopTracerProvider())
}
