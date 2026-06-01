package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/YASSERRMD/sql-sage/backend/internal/config"
	"github.com/YASSERRMD/sql-sage/backend/internal/database"
	"github.com/YASSERRMD/sql-sage/backend/internal/models"
	"github.com/YASSERRMD/sql-sage/backend/pkg/crypto"
	"github.com/YASSERRMD/sql-sage/backend/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}
	log := logger.New(cfg.LogLevel)

	db, err := database.Open(cfg, log)
	if err != nil {
		log.Error("db open", "err", err)
		os.Exit(1)
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Error("automigrate", "err", err)
		os.Exit(1)
	}

	if err := seed(db, cfg, log); err != nil {
		log.Error("seed", "err", err)
		os.Exit(1)
	}
	log.Info("seed complete")
}

func seed(db *gorm.DB, cfg *config.Config, log interface{}) error {
	email := os.Getenv("SEED_ADMIN_EMAIL")
	password := os.Getenv("SEED_ADMIN_PASSWORD")
	name := os.Getenv("SEED_ADMIN_NAME")
	if email == "" || password == "" {
		return errors.New("SEED_ADMIN_EMAIL and SEED_ADMIN_PASSWORD required")
	}
	if name == "" {
		name = "Administrator"
	}

	var existing models.User
	tx := db.WithContext(context.Background()).Where("email = ?", email).First(&existing)
	if tx.Error == nil {
		log_iface(log).Info("admin already exists", "email", email)
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin := &models.User{
		Email:        email,
		Name:         name,
		PasswordHash: string(hash),
		Role:         models.RoleAdmin,
		IsActive:     true,
	}
	if err := db.Create(admin).Error; err != nil {
		return err
	}

	if _, err := crypto.New(cfg.EncryptionKey); err != nil {
		return fmt.Errorf("encryption key invalid: %w", err)
	}

	log_iface(log).Info("admin created", "email", email, "createdAt", admin.CreatedAt.Format(time.RFC3339))
	return nil
}

type slogger interface {
	Info(msg string, args ...any)
}

func log_iface(l interface{}) slogger {
	if s, ok := l.(slogger); ok {
		return s
	}
	return nopLogger{}
}

type nopLogger struct{}

func (nopLogger) Info(string, ...any) {}
