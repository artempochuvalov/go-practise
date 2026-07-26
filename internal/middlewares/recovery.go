package middlewares

import (
	"log/slog"
	"net/http"
)

func RecoverMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					log.Error(
						"panic recovered",
						"method", req.Method,
						"path", req.URL.Path,
						"panic", recovered,
					)

					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, req)
		})
	}
}
