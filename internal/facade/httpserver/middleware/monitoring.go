package middleware

import (
	"net/http"
)

func (m *middleware) MonitoringMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/metrics" {
			h = m.mon.WrapHandler(r.URL.Path, h)
		}
		h.ServeHTTP(w, r)
	})
}
