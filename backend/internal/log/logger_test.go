package log

import "testing"

func TestNewLogger(t *testing.T) {
	logger := New("development")
	if logger == nil {
		t.Fatalf("expected logger instance")
	}

	prodLogger := New("production")
	if prodLogger == nil {
		t.Fatalf("expected logger instance")
	}
}
