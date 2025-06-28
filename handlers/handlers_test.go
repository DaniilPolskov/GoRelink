package handlers

import (
	"GoRelink/storage"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := ShortenHandler(store)

	t.Run("valid URL", func(t *testing.T) {
		body := map[string]string{"url": "go.dev"}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "shortURL") {
			t.Fatalf("expected response to contain 'shortURL', got %s", rr.Body.String())
		}
	})

	t.Run("invalid method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/shorten", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("expected 405, got %d", rr.Code)
		}
	})
}

func TestBatchShortenHandler(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := BatchShortenHandler(store)

	t.Run("valid batch request", func(t *testing.T) {
		body := map[string]interface{}{
			"url":   "go.dev",
			"count": 3,
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "shortURLs") {
			t.Fatalf("expected 'shortURLs' in response, got %s", rr.Body.String())
		}
	})

	t.Run("invalid count", func(t *testing.T) {
		body := map[string]interface{}{
			"url":   "go.dev",
			"count": 11,
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", rr.Code)
		}
	})
}

func TestRedirectHandler(t *testing.T) {
	store := storage.NewMemoryStore()
	store.Save("abc123", "https://go.dev")
	handler := RedirectHandler(store)

	t.Run("existing ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/gorelink/abc123", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusFound {
			t.Fatalf("expected 302 Found, got %d", rr.Code)
		}
		loc := rr.Header().Get("Location")
		if loc != "https://go.dev" {
			t.Fatalf("expected redirect to https://go.dev, got %s", loc)
		}
	})

	t.Run("nonexistent ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/gorelink/unknown", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("expected 404 Not Found, got %d", rr.Code)
		}
	})
}
