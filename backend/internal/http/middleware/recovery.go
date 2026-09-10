package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"city-of-lies/backend/internal/http/respond"
)

func Recovery(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := string(debug.Stack())
					reqID, _ := r.Context().Value(RequestIDKey).(string)
					logger.Error("panic_recovered",
						slog.Any("panic", rec),
						slog.String("stack", stack),
						slog.String("request_id", reqID),
					)
					respond.Error(w, http.StatusInternalServerError, "Internal Server Error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
