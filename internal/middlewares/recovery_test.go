package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoveryMiddleware(t *testing.T) {
	log, buf := newTestLogger()

	handler := RecoverMiddleware(log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest(http.MethodPost, "/todos", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rr.Code,
		)
	}

	if !strings.Contains(rr.Body.String(), "internal server error") {
		t.Fatalf("unexpected response body: %s", rr.Body.String())
	}

	loggedString := buf.String()

	if !strings.Contains(loggedString, "panic recovered") {
		t.Fatalf("invalid log info: %s", loggedString)
	}

	if !strings.Contains(loggedString, fmt.Sprintf("method=%s", http.MethodPost)) {
		t.Fatalf("log should have contained method %s", http.MethodPost)
	}

	if !strings.Contains(loggedString, "path=/todos") {
		t.Fatalf("log should have contained path %s", "/todos")
	}

	if !strings.Contains(loggedString, `panic="test panic"`) {
		t.Fatalf("log should have panic context: %s", loggedString)
	}
}

func TestRecoveryMiddleware_NoPanic(t *testing.T) {
	log, _ := newTestLogger()

	called := false

	handler := RecoverMiddleware(log)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusNoContent)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Fatal("next handler should have been called")
	}

	if rr.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNoContent,
			rr.Code,
		)
	}
}
