package httpx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setup() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAbortWithError(t *testing.T) {
	r := setup()
	r.GET("/x", func(c *gin.Context) {
		AbortWithError(c, http.StatusBadRequest, "BAD", "msg", map[string]any{"a": 1})
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/x", nil)
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var body ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Error.Code != "BAD" {
		t.Fatalf("expected BAD, got %s", body.Error.Code)
	}
}

func TestOKAndCreated(t *testing.T) {
	r := setup()
	r.GET("/ok", func(c *gin.Context) { OK(c, gin.H{"x": 1}) })
	r.GET("/c", func(c *gin.Context) { Created(c, gin.H{"y": 2}) })
	r.GET("/n", func(c *gin.Context) { NoContent(c) })

	for _, tc := range []struct {
		path string
		code int
	}{
		{"/ok", 200}, {"/c", 201}, {"/n", 204},
	} {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", tc.path, nil)
		r.ServeHTTP(w, req)
		if w.Code != tc.code {
			t.Fatalf("path %s: expected %d, got %d", tc.path, tc.code, w.Code)
		}
	}
}
