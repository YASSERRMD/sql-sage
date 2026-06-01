package repositories

import (
	"context"
	"errors"

	"github.com/YASSERRMD/sql-sage/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, t *models.RefreshToken) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *RefreshTokenRepository) FindByHash(ctx context.Context, hash string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	if err := r.db.WithContext(ctx).Where("token_hash = ?", hash).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).Where("id = ?", id).Update("revoked", true).Error
}

func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).Where("user_id = ?", userID).Update("revoked", true).Error
}
