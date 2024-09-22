package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"runtime/debug"
)

func (m *middleware) RecoverMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				stacktrace := string(debug.Stack())
				switch t := r.(type) {
				case string:
					err = fmt.Errorf(`panic: %s, stacktrace: %s`, t, stacktrace)
				case error:
					err = fmt.Errorf(`panic: %w, stacktrace: %s`, t, stacktrace)
				default:
					err = errors.New("unknown panic")
				}
				m.lg.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
