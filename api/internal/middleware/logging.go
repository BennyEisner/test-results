package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs details about each incoming HTTP request.
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the response writer to capture status code
			rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rr, r)

			duration := time.Since(start)
			logger.Info("HTTP request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rr.statusCode),
				slog.String("remote", r.RemoteAddr),
				slog.Duration("duration", duration),
			)
		})
	}
}

// responseRecorder wraps http.ResponseWriter to record the status code.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

