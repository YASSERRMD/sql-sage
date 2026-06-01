package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiterBlocks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rl := NewRateLimiter(2)
	r := gin.New()
	r.GET("/x", rl.Middleware(), func(c *gin.Context) { c.Status(204) })

	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/x", nil)
		req.RemoteAddr = "1.2.3.4:5555"
		r.ServeHTTP(w, req)
		if w.Code != 204 {
			t.Fatalf("call %d: expected 204, got %d", i, w.Code)
		}
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/x", nil)
	req.RemoteAddr = "1.2.3.4:5555"
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w.Code)
	}
}

func TestRateLimiterResetsWindow(t *testing.T) {
	rl := NewRateLimiter(1)
	rl.window = 10 * time.Millisecond
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", rl.Middleware(), func(c *gin.Context) { c.Status(204) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/x", nil)
	req.RemoteAddr = "5.6.7.8:1234"
	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	time.Sleep(20 * time.Millisecond)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/x", nil)
	req2.RemoteAddr = "5.6.7.8:1234"
	r.ServeHTTP(w2, req2)
	if w2.Code != 204 {
		t.Fatalf("expected 204 after reset, got %d", w2.Code)
	}
}
