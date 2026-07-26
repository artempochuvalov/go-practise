package middlewares

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newLogger() (*slog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer

	log := slog.New(
		slog.NewTextHandler(&buf, nil),
	)

	return log, &buf
}

func TestLoggingMiddleware(t *testing.T) {
	log, buf := newLogger()

	handler := LoggingMiddleware(log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/todos", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			rr.Code,
		)
	}

	loggedString := buf.String()

	if !strings.Contains(loggedString, fmt.Sprintf("status=%d", http.StatusCreated)) {
		t.Fatalf("log should have contained status %d", http.StatusCreated)
	}

	if !strings.Contains(loggedString, fmt.Sprintf("method=%s", http.MethodPost)) {
		t.Fatalf("log should have contained method %s", http.MethodPost)
	}

	if !strings.Contains(loggedString, "path=/todos") {
		t.Fatalf("log should have contained path %s", "/todos")
	}
}

func TestLoggingMiddleware_DefaultStatus(t *testing.T) {
	log, buf := newLogger()

	var called bool
	handler := LoggingMiddleware(log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest(http.MethodPost, "/todos", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rr.Code,
		)
	}

	if !called {
		t.Fatalf("handler should have been called")
	}

	loggedString := buf.String()

	if !strings.Contains(loggedString, fmt.Sprintf("status=%d", http.StatusOK)) {
		t.Fatalf("log should have contained status %d", http.StatusOK)
	}
}
