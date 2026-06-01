package database

import (
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
)

func TestRunMigrationsSkipWithoutDB(t *testing.T) {
	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("set RUN_DB_TESTS=1 to run with a real database")
	}
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Fatal("TEST_DB_DSN required")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	path, _ := filepath.Abs("./../../migrations")
	if err := runMigrationsWith(db, path, slog.Default()); err != nil {
		t.Fatal(err)
	}
}
