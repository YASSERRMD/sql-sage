package database

import (
	"fmt"

	"github.com/YASSERRMD/sql-sage/backend/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Provider{},
		&models.Analysis{},
		&models.RefreshToken{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}
