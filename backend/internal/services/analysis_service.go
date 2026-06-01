package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/YASSERRMD/sql-sage/backend/internal/analysis"
	"github.com/YASSERRMD/sql-sage/backend/internal/models"
	"github.com/YASSERRMD/sql-sage/backend/internal/repositories"
	"github.com/google/uuid"
)

var ErrAnalysisNotFound = errors.New("analysis not found")

type AnalysisService struct {
	engine      *analysis.Engine
	repo        *repositories.AnalysisRepository
	providerSvc *ProviderService
}

func NewAnalysisService(e *analysis.Engine, r *repositories.AnalysisRepository, p *ProviderService) *AnalysisService {
	return &AnalysisService{engine: e, repo: r, providerSvc: p}
}

type CreateAnalysisInput struct {
	ObjectName string             `json:"objectName"`
	ObjectType models.ObjectType  `json:"objectType"`
	SourceCode string             `json:"sourceCode"`
	ProviderID *uuid.UUID         `json:"providerId,omitempty"`
}

func (s *AnalysisService) Create(ctx context.Context, userID uuid.UUID, in CreateAnalysisInput) (*models.Analysis, error) {
	if in.SourceCode == "" {
		return nil, fmt.Errorf("%w: sourceCode required", ErrValidation)
	}
	if !in.ObjectType.Valid() {
		in.ObjectType = models.ObjectUnknown
	}
	res, err := s.engine.Analyze(ctx, analysis.AnalyzeInput{
		UserID:     userID,
		ProviderID: in.ProviderID,
		ObjectName: in.ObjectName,
		ObjectType: in.ObjectType,
		SourceCode: in.SourceCode,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, userID, res.AnalysisID)
}

func (s *AnalysisService) Get(ctx context.Context, userID, id uuid.UUID) (*models.Analysis, error) {
	a, err := s.repo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, ErrAnalysisNotFound
	}
	return a, nil
}

func (s *AnalysisService) Delete(ctx context.Context, userID, id uuid.UUID) error {
	return s.repo.Delete(ctx, userID, id)
}

type AnalysisListItem struct {
	ID         uuid.UUID `json:"id"`
	ObjectName string    `json:"objectName"`
	ObjectType string    `json:"objectType"`
	Summary    string    `json:"summary"`
	RiskScore  string    `json:"riskScore"`
	CreatedAt  string    `json:"createdAt"`
}

func (s *AnalysisService) List(ctx context.Context, f repositories.ListFilter) ([]AnalysisListItem, int64, error) {
	items, total, err := s.repo.List(ctx, f)
	if err != nil {
		return nil, 0, err
	}
	out := make([]AnalysisListItem, 0, len(items))
	for _, a := range items {
		out = append(out, AnalysisListItem{
			ID:         a.ID,
			ObjectName: a.ObjectName,
			ObjectType: string(a.ObjectType),
			Summary:    a.Summary,
			RiskScore:  a.RiskScore,
			CreatedAt:  a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return out, total, nil
}

func (s *AnalysisService) DecodeResult(a *models.Analysis) (*analysis.Schema, error) {
	var s2 analysis.Schema
	if err := json.Unmarshal([]byte(a.Result), &s2); err != nil {
		return nil, err
	}
	return &s2, nil
}
