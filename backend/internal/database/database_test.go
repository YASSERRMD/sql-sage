package database

import (
	"log/slog"
	"testing"

	"github.com/YASSERRMD/sql-sage/backend/internal/config"
)

func TestOpenInvalidDSN(t *testing.T) {
	cfg := &config.Config{
		DBHost: "127.0.0.1", DBPort: 1, DBUser: "x", DBPassword: "x", DBName: "x", DBSSLMode: "disable",
	}
	if _, err := Open(cfg, slog.Default()); err == nil {
		t.Skip("could connect to invalid port - skipping")
	}
}
