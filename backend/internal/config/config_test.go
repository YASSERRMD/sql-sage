package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	os.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef")
	os.Setenv("ENCRYPTION_KEY", "01234567890123456789012345678901")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("ENCRYPTION_KEY")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.HTTPPort != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.HTTPPort)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("expected localhost, got %s", cfg.DBHost)
	}
}

func TestLoadValidation(t *testing.T) {
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("ENCRYPTION_KEY")
	if _, err := Load(); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestDSN(t *testing.T) {
	cfg := &Config{
		DBHost: "h", DBPort: 5432, DBUser: "u", DBPassword: "p", DBName: "n", DBSSLMode: "disable",
	}
	dsn := cfg.DSN()
	if dsn == "" {
		t.Fatal("expected non-empty DSN")
	}
}

func TestSplitCSV(t *testing.T) {
	out := splitCSV("a, b , c,")
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	if out[0] != "a" || out[1] != "b" || out[2] != "c" {
		t.Fatalf("unexpected: %v", out)
	}
}
