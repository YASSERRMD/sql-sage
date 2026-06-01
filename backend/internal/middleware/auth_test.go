package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YASSERRMD/sql-sage/backend/internal/auth"
	"github.com/YASSERRMD/sql-sage/backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func newSvc(t *testing.T) *auth.Service {
	t.Helper()
	gin.SetMode(gin.TestMode)
	c := &config.Config{
		JWTSecret:     "0123456789abcdef0123456789abcdef",
		EncryptionKey: "01234567890123456789012345678901",
		JWTAccessTTL:  900 * 1_000_000_000,
		JWTRefreshTTL: 86400 * 1_000_000_000,
	}
	return auth.NewService(c)
}

func TestAuthRequiredMissing(t *testing.T) {
	svc := newSvc(t)
	r := gin.New()
	r.GET("/x", AuthRequired(svc), func(c *gin.Context) { c.Status(204) })
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/x", nil)
	r.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequiredValid(t *testing.T) {
	svc := newSvc(t)
	uid := uuid.New()
	pair, err := svc.IssueTokens(uid, "user", "u@x.com")
	if err != nil {
		t.Fatal(err)
	}
	r := gin.New()
	r.GET("/x", AuthRequired(svc), func(c *gin.Context) {
		if GetUserID(c) != uid.String() {
			t.Fatalf("uid mismatch: %s", GetUserID(c))
		}
		c.Status(204)
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/x", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Fatalf("expected 204, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestRequireRole(t *testing.T) {
	r := gin.New()
	r.GET("/admin", AuthRequired(newSvc(t)), RequireRole("admin"), func(c *gin.Context) { c.Status(204) })

	svc := newSvc(t)
	uid := uuid.New()
	userPair, _ := svc.IssueTokens(uid, "user", "u@x.com")
	adminPair, _ := svc.IssueTokens(uid, "admin", "a@x.com")

	hit := func(tok string) int {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+tok)
		r.ServeHTTP(w, req)
		return w.Code
	}
	if got := hit(userPair.AccessToken); got != 403 {
		var body map[string]any
		_ = json.Unmarshal(nil, &body)
		t.Fatalf("expected 403, got %d", got)
	}
	if got := hit(adminPair.AccessToken); got != 204 {
		t.Fatalf("expected 204, got %d", got)
	}
}
