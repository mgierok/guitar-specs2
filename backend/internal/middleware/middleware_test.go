package middleware

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestID(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := RequestIDFromContext(r.Context())
		if requestID == "" {
			t.Fatalf("expected request_id in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("X-Request-Id", "test-id")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Header().Get("X-Request-Id") != "test-id" {
		t.Fatalf("expected request id header to propagate")
	}
}

func TestRequestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{}))

	handler := RequestID(RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("ok"))
	})))

	req := httptest.NewRequest(http.MethodGet, "/guitars", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	output := buf.String()
	if !strings.Contains(output, "request_id=") {
		t.Fatalf("expected request_id in logs")
	}
	if !strings.Contains(output, "status=201") {
		t.Fatalf("expected status in logs")
	}
}

func TestRequestIDFromContextEmpty(t *testing.T) {
	if RequestIDFromContext(context.Background()) != "" {
		t.Fatalf("expected empty request_id for empty context")
	}
}
