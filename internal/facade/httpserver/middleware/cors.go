package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func (m *middleware) CorsMiddleware(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		Debug:            false,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
	})

	return c.Handler(h)
}
