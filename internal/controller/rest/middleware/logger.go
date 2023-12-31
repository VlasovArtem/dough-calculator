package middleware

import (
	"net/http"

	"github.com/rs/zerolog"
)

func LoggerMiddleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Msg("request started")

			next.ServeHTTP(w, r)

			logger.Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Msg("request completed")
		})
	}
}
