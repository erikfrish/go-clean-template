package middleware

import (
	"net/http"
	"unicode/utf8"
)

func (m *middleware) ValidationMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !utf8.ValidString(r.URL.Path) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.ServeHTTP(w, r)
	})
}
