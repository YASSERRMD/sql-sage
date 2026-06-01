package auth

import (
	"os"
	"testing"
	"time"

	"github.com/YASSERRMD/sql-sage/backend/internal/config"
	"github.com/google/uuid"
)

func cfg() *config.Config {
	os.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef")
	os.Setenv("ENCRYPTION_KEY", "01234567890123456789012345678901")
	c, _ := config.Load()
	return c
}

func TestIssueAndParse(t *testing.T) {
	s := NewService(cfg())
	uid := uuid.New()
	pair, err := s.IssueTokens(uid, "user", "u@x.com")
	if err != nil {
		t.Fatal(err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatal("expected non-empty tokens")
	}
	claims, err := s.ParseAccess(pair.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != uid.String() {
		t.Fatalf("uid mismatch: %s", claims.UserID)
	}
	if claims.Role != "user" {
		t.Fatalf("role mismatch")
	}
	if claims.Type != "access" {
		t.Fatalf("type mismatch")
	}
}

func TestParseInvalid(t *testing.T) {
	s := NewService(cfg())
	if _, err := s.ParseAccess("not-a-token"); err == nil {
		t.Fatal("expected error")
	}
}

func TestHashAndTTL(t *testing.T) {
	s := NewService(cfg())
	tok, _ := s.RandomToken()
	h := s.HashToken(tok)
	if h == "" || len(h) != 64 {
		t.Fatalf("expected 64 char hex, got %d", len(h))
	}
	if s.RefreshTTL() <= 0 {
		t.Fatal("expected positive refresh TTL")
	}
	_ = time.Now()
}
