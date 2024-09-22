package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type ContextKey string

const contextKeyRequestID ContextKey = "requestID"

var ignorePaths = []string{ //nolint:gochecknoglobals //calls from infra should not be logged
	"/metrics",
	"/api/live",
	"/api/ready",
}

type ResponseLogger struct {
	w          http.ResponseWriter
	statusCode int
	body       []string
}

func (r *ResponseLogger) Header() http.Header {
	return r.w.Header()
}

func (r *ResponseLogger) Write(res []byte) (int, error) {
	r.body = append(r.body, string(res))
	return r.w.Write(res)
}

func (r *ResponseLogger) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.w.WriteHeader(statusCode)
}

func (r *ResponseLogger) fullResponse() string {
	if len(r.body) == 1 {
		return r.body[0]
	}
	var res string
	for i, b := range r.body {
		res += fmt.Sprintf("written %d: %s\n", i+1, b)
	}
	return res
}

func (m *middleware) RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New()
		ctx := context.WithValue(r.Context(), contextKeyRequestID, reqID.String())
		r = r.WithContext(ctx)

		body, err := readReqBody(r)
		if err != nil {
			m.lg.Error(reqID.String(), err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rl := &ResponseLogger{w: w, body: make([]string, 0)}
		if !slices.Contains(ignorePaths, r.RequestURI) {
			if strings.Contains(r.RequestURI, "upload") {
				m.lg.Info(reqID.String(), fmt.Sprintf("Request: %s %s %s %d", r.Method, r.RequestURI, r.RemoteAddr, len(body)))
			} else {
				m.lg.Info(reqID.String(), fmt.Sprintf("Request: %s %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, string(body)))
			}
			defer func() {
				if strings.Contains(r.RequestURI, "search") || strings.Contains(r.RequestURI, "list") {
					m.lg.Info(reqID.String(), fmt.Sprintf("Response: %s %s %d %d",
						r.Method, r.RequestURI, rl.statusCode, len([]byte(rl.fullResponse()))))
				} else {
					m.lg.Info(reqID.String(), fmt.Sprintf("Response: %s %s %d %s",
						r.Method, r.RequestURI, rl.statusCode, rl.fullResponse()))
				}
			}()
		}

		h.ServeHTTP(rl, r)
	})
}

func readReqBody(r *http.Request) ([]byte, error) {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(r.Body); err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(buf)

	return buf.Bytes(), nil
}
