package rewire

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewHTTPTransport wraps base with OTel trace propagation so outgoing requests
// are recorded as child spans of the current trace. If base is nil,
// http.DefaultTransport is used.
func NewHTTPTransport(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return otelhttp.NewTransport(base)
}

// NewHTTPMiddleware wraps h so that incoming requests extract trace context
// from headers and create a server span for each request.
func NewHTTPMiddleware(h http.Handler) http.Handler {
	return otelhttp.NewHandler(h, "http.server")
}
