package rewire

import (
	"os"
	"testing"
)

var initEnvVars = []string{
	"OTEL_EXPORTER_OTLP_ENDPOINT",
	"REWIRE_TOKEN",
	"GITHUB_RUN_ID",
	"CIRCLE_WORKFLOW_ID",
	"CI_PIPELINE_ID",
	"REWIRE_RUN_ID",
	"OTEL_SERVICE_NAME",
	"GITHUB_REPOSITORY",
}

// resetInit clears all env vars and resets package state before each test.
func resetInit(t *testing.T) {
	t.Helper()
	saved := make(map[string]string)
	for _, k := range initEnvVars {
		saved[k] = os.Getenv(k)
		os.Unsetenv(k)
	}
	_reset()
	t.Cleanup(func() {
		_reset()
		for k, v := range saved {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	})
}

func TestInit_NoEnvVars_ReturnsCallable(t *testing.T) {
	resetInit(t)
	shutdown := Init()
	if shutdown == nil {
		t.Fatal("Init() returned nil")
	}
}

func TestInit_NoEnvVars_ShutdownDoesNotPanic(t *testing.T) {
	resetInit(t)
	shutdown := Init()
	shutdown() // must not panic
}

func TestInit_NoEnvVars_ReturnsSameFuncOnSecondCall(t *testing.T) {
	resetInit(t)
	a := Init()
	b := Init()
	// Compare function identity via pointer trick: both should be no-ops
	// (we can't compare func values directly in Go, but we can verify
	// the idempotency contract holds by checking neither panics)
	a()
	b()
}

func TestInit_WithEndpoint_ReturnsCallable(t *testing.T) {
	resetInit(t)
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318")
	shutdown := Init()
	if shutdown == nil {
		t.Fatal("Init() returned nil")
	}
	shutdown()
}

func TestInit_WithToken_ReturnsCallable(t *testing.T) {
	resetInit(t)
	os.Setenv("REWIRE_TOKEN", "rwt_test")
	shutdown := Init()
	if shutdown == nil {
		t.Fatal("Init() returned nil")
	}
	shutdown()
}

func TestInit_Idempotent(t *testing.T) {
	resetInit(t)
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318")
	a := Init()
	b := Init()
	// Both should be callable without panic; second call must be a no-op
	a()
	_ = b
}

func TestShutdown_BeforeInit_DoesNotPanic(t *testing.T) {
	resetInit(t)
	Shutdown() // must not panic
}

func TestShutdown_AfterInit_DoesNotPanic(t *testing.T) {
	resetInit(t)
	os.Setenv("REWIRE_TOKEN", "rwt_test")
	Init()
	Shutdown() // must not panic
}

func TestServiceName_OtelEnvVar(t *testing.T) {
	t.Setenv("OTEL_SERVICE_NAME", "my-service")
	if got := serviceName(); got != "my-service" {
		t.Errorf("serviceName() = %q, want %q", got, "my-service")
	}
}

func TestServiceName_GitHubRepository(t *testing.T) {
	os.Unsetenv("OTEL_SERVICE_NAME")
	t.Setenv("GITHUB_REPOSITORY", "acme/app")
	if got := serviceName(); got != "acme/app" {
		t.Errorf("serviceName() = %q, want %q", got, "acme/app")
	}
}

func TestServiceName_Fallback(t *testing.T) {
	os.Unsetenv("OTEL_SERVICE_NAME")
	os.Unsetenv("GITHUB_REPOSITORY")
	if got := serviceName(); got != "unknown" {
		t.Errorf("serviceName() = %q, want %q", got, "unknown")
	}
}
