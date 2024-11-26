package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JohnGrimm/rate-limiter-fg/internal/middleware"
	"github.com/JohnGrimm/rate-limiter-fg/internal/services"
	"github.com/JohnGrimm/rate-limiter-fg/internal/storage"
)

func TestRateLimiterMiddlewareByIP(t *testing.T) {
	storage := &storage.MockStorage{
		Data: make(map[string]interface{}),
	}

	limit := *services.NewLimiter(storage, 5, 5, 10)
	limit.ProcessKeysFromFile()
	middleware := middleware.NewRateLimiterMiddleware(limit)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1"

	for i := 0; i < 5; i++ {
		rr := httptest.NewRecorder()
		middleware.Middleware(handler).ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rr.Code)
		}
	}

	rr := httptest.NewRecorder()
	middleware.Middleware(handler).ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status Too Many Requests, got %v", rr.Code)
	}
}

func TestRateLimiterMiddlewareByAPIKey(t *testing.T) {
	storage := &storage.MockStorage{
		Data: make(map[string]interface{}),
	}

	limit := *services.NewLimiter(storage, 10, 10, 10)
	limit.ProcessKeysFromFile()
	middleware := middleware.NewRateLimiterMiddleware(limit)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "03386fab-cf28-44e5-a99f-19bdc866ce73")

	for i := 0; i < 15; i++ {
		rr := httptest.NewRecorder()
		middleware.Middleware(handler).ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rr.Code)
		}
	}

	rr := httptest.NewRecorder()
	middleware.Middleware(handler).ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status Too Many Requests, got %v", rr.Code)
	}
}

func TestRateLimiterMiddlewareByAPIKeyNotFound(t *testing.T) {
	storage := &storage.MockStorage{Data: make(map[string]interface{})}

	limit := *services.NewLimiter(storage, 5, 5, 5)

	limit.ProcessKeysFromFile()
	middleware := middleware.NewRateLimiterMiddleware(limit)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "notfound")
	req.RemoteAddr = "127.0.0.1"

	for i := 0; i < 5; i++ {
		rr := httptest.NewRecorder()
		middleware.Middleware(handler).ServeHTTP(rr, req)
		if rr.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status Too Many Requests, got %v", rr.Code)
		}
	}

	rr := httptest.NewRecorder()
	middleware.Middleware(handler).ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status Too Many Requests, got %v", rr.Code)
	}
}

func TestRateLimiterIfBlockTheUserWhenTheLimitIsReached(t *testing.T) {
	storage := &storage.MockStorage{Data: make(map[string]interface{})}

	limit := *services.NewLimiter(storage, 10, 10, 10)

	limit.ProcessKeysFromFile()
	middleware := middleware.NewRateLimiterMiddleware(limit)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "03386fab-cf28-44e5-a99f-19bdc866ce73")

	for i := 0; i < 15; i++ {
		rr := httptest.NewRecorder()
		middleware.Middleware(handler).ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rr.Code)
		}
	}

	rr := httptest.NewRecorder()
	middleware.Middleware(handler).ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status Too Many Requests, got %v", rr.Code)
	}

}
