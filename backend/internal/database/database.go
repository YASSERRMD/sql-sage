package database

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/YASSERRMD/sql-sage/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(cfg *config.Config, slogger *slog.Logger) (*gorm.DB, error) {
	gormLogger := logger.New(
		getGormWriter(slogger),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	return db, nil
}

func getGormWriter(slogger *slog.Logger) logger.Writer {
	return &slogWriter{slogger: slogger}
}

type slogWriter struct {
	slogger *slog.Logger
}

func (w *slogWriter) Printf(format string, args ...any) {
	if w.slogger == nil {
		fmt.Fprintf(os.Stderr, format, args...)
		return
	}
	w.slogger.Info(fmt.Sprintf(format, args...))
}
