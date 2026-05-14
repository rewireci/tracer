package rewire

import (
	"net/http"
	"testing"
)

func TestNewHTTPTransport_NilBase_UsesDefault(t *testing.T) {
	rt := NewHTTPTransport(nil)
	if rt == nil {
		t.Fatal("NewHTTPTransport(nil) returned nil")
	}
}

func TestNewHTTPTransport_CustomBase_Wraps(t *testing.T) {
	base := http.DefaultTransport
	rt := NewHTTPTransport(base)
	if rt == nil {
		t.Fatal("NewHTTPTransport(base) returned nil")
	}
	// Verify it's a different type (wrapped), not the same object
	if rt == base {
		t.Error("NewHTTPTransport should wrap the base transport, not return it as-is")
	}
}

func TestNewHTTPMiddleware_WrapsHandler(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	h := NewHTTPMiddleware(inner)
	if h == nil {
		t.Fatal("NewHTTPMiddleware returned nil")
	}
}
