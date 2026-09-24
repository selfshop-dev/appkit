package appkit

import (
	"net/http"
	"sync/atomic"
)

// deferredHandler is an HTTP handler whose implementation can be replaced
// atomically after creation.
type deferredHandler struct {
	handler atomic.Pointer[http.Handler]
}

// NewDeferredHandler creates an HTTP handler whose implementation can be
// replaced with Set.
//
// If base is nil, the handler initially responds with HTTP 503 Service
// Unavailable. This allows an endpoint to be registered before the underlying
// application or service is ready.
//
// The returned handler can be used anywhere an [http.Handler] is expected.
// Calling Set replaces the current handler atomically.
func NewDeferredHandler(base http.HandlerFunc) *deferredHandler {
	h := &deferredHandler{}
	v := http.Handler(base)
	if base == nil {
		f := http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			http.Error(
				w,
				http.StatusText(http.StatusServiceUnavailable),
				http.StatusServiceUnavailable,
			)
		})
		v = http.Handler(f)
	}
	h.handler.Store(&v)
	return h
}

// Set atomically replaces the current HTTP handler.
//
// The new handler is used by subsequent ServeHTTP calls.
func (h *deferredHandler) Set(next http.Handler) {
	h.handler.Store(&next)
}

// ServeHTTP implements [http.Handler].
//
// The currently configured handler is loaded atomically and used to serve the
// request.
func (h *deferredHandler) ServeHTTP(
	w http.ResponseWriter, r *http.Request,
) {
	(*h.handler.Load()).ServeHTTP(w, r)
}
