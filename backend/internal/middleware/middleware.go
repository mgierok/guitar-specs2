package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"log/slog"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = newRequestID()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		r = r.WithContext(ctx)
		w.Header().Set("X-Request-Id", requestID)
		next.ServeHTTP(w, r)
	})
}

func RequestIDFromContext(ctx context.Context) string {
	value := ctx.Value(requestIDKey)
	requestID, _ := value.(string)
	return requestID
}

func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			reqID := RequestIDFromContext(r.Context())

			wrapped := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logger.Info("request",
				slog.String("request_id", reqID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", wrapped.status),
				slog.Int64("duration_ms", duration.Milliseconds()),
			)
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func newRequestID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(buf)
}
