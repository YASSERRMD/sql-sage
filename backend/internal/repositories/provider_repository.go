package repositories

import (
	"context"
	"errors"

	"github.com/YASSERRMD/sql-sage/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProviderRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) *ProviderRepository { return &ProviderRepository{db: db} }

func (r *ProviderRepository) Create(ctx context.Context, p *models.Provider) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProviderRepository) Update(ctx context.Context, p *models.Provider) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *ProviderRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Provider{}).Error
}

func (r *ProviderRepository) FindByID(ctx context.Context, userID, id uuid.UUID) (*models.Provider, error) {
	var p models.Provider
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProviderRepository) Get(ctx context.Context, userID, id uuid.UUID) (*models.Provider, error) {
	return r.FindByID(ctx, userID, id)
}

func (r *ProviderRepository) List(ctx context.Context, userID uuid.UUID) ([]models.Provider, error) {
	var ps []models.Provider
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&ps).Error; err != nil {
		return nil, err
	}
	return ps, nil
}

func (r *ProviderRepository) FindDefault(ctx context.Context, userID uuid.UUID) (*models.Provider, error) {
	var p models.Provider
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_default = ?", userID, true).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProviderRepository) ClearDefault(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Provider{}).Where("user_id = ?", userID).Update("is_default", false).Error
}

func (r *ProviderRepository) SetDefault(ctx context.Context, userID, id uuid.UUID) error {
	if err := r.ClearDefault(ctx, userID); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&models.Provider{}).Where("id = ? AND user_id = ?", id, userID).Update("is_default", true).Error
}
