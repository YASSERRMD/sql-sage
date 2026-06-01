package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/YASSERRMD/sql-sage/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnalysisRepository struct {
	db *gorm.DB
}

func NewAnalysisRepository(db *gorm.DB) *AnalysisRepository { return &AnalysisRepository{db: db} }

func (r *AnalysisRepository) Create(ctx context.Context, a *models.Analysis) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *AnalysisRepository) FindByID(ctx context.Context, userID, id uuid.UUID) (*models.Analysis, error) {
	var a models.Analysis
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *AnalysisRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Analysis{}).Error
}

type ListFilter struct {
	UserID     uuid.UUID
	Query      string
	ObjectType string
	Risk       string
	Limit      int
	Offset     int
}

func (r *AnalysisRepository) List(ctx context.Context, f ListFilter) ([]models.Analysis, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Analysis{}).Where("user_id = ?", f.UserID)
	if s := strings.TrimSpace(f.Query); s != "" {
		like := "%" + s + "%"
		q = q.Where("object_name ILIKE ? OR summary ILIKE ?", like, like)
	}
	if f.ObjectType != "" {
		q = q.Where("object_type = ?", f.ObjectType)
	}
	if f.Risk != "" {
		q = q.Where("risk_score = ?", f.Risk)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []models.Analysis
	if err := q.Order("created_at DESC").Limit(f.Limit).Offset(f.Offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *AnalysisRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Analysis{}).Where("user_id = ?", userID).Count(&n).Error
	return n, err
}

func (r *AnalysisRepository) CountByType(ctx context.Context, userID uuid.UUID, t string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Analysis{}).Where("user_id = ? AND object_type = ?", userID, t).Count(&n).Error
	return n, err
}

func (r *AnalysisRepository) CountByRisk(ctx context.Context, userID uuid.UUID, risk string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Analysis{}).Where("user_id = ? AND risk_score = ?", userID, risk).Count(&n).Error
	return n, err
}
