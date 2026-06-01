package services

import (
	"context"
	"time"

	"github.com/YASSERRMD/sql-sage/backend/internal/models"
	"github.com/YASSERRMD/sql-sage/backend/internal/repositories"
	"github.com/google/uuid"
)

type DashboardService struct {
	analyses  *repositories.AnalysisRepository
	providers *repositories.ProviderRepository
}

func NewDashboardService(a *repositories.AnalysisRepository, p *repositories.ProviderRepository) *DashboardService {
	return &DashboardService{analyses: a, providers: p}
}

type Summary struct {
	TotalAnalyses   int64                `json:"totalAnalyses"`
	TotalProcedures int64                `json:"totalProcedures"`
	TotalFunctions  int64                `json:"totalFunctions"`
	TotalPackages   int64                `json:"totalPackages"`
	HighRisk        int64                `json:"highRiskFindings"`
	ProviderUsage   []ProviderUsageEntry `json:"providerUsage"`
}

type ProviderUsageEntry struct {
	ProviderID   string `json:"providerId"`
	ProviderName string `json:"providerName"`
	Count        int64  `json:"count"`
}

func (s *DashboardService) Summary(ctx context.Context, userID uuid.UUID) (*Summary, error) {
	total, err := s.analyses.CountByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	proc, _ := s.analyses.CountByType(ctx, userID, string(models.ObjectProcedure))
	fun, _ := s.analyses.CountByType(ctx, userID, string(models.ObjectFunction))
	pkg, _ := s.analyses.CountByType(ctx, userID, string(models.ObjectPackage))
	high, _ := s.analyses.CountByRisk(ctx, userID, string(models.RiskHigh))
	critical, _ := s.analyses.CountByRisk(ctx, userID, string(models.RiskCritical))

	providers, err := s.providers.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	usage := make([]ProviderUsageEntry, 0, len(providers))
	for _, p := range providers {
		c, _ := s.analyses.CountByUser(ctx, userID) // fallback
		_ = c
		usage = append(usage, ProviderUsageEntry{
			ProviderID:   p.ID.String(),
			ProviderName: p.Name,
			Count:        0,
		})
	}
	return &Summary{
		TotalAnalyses:   total,
		TotalProcedures: proc,
		TotalFunctions:  fun,
		TotalPackages:   pkg,
		HighRisk:        high + critical,
		ProviderUsage:   usage,
	}, nil
}

type TrendPoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

func (s *DashboardService) Trend(ctx context.Context, userID uuid.UUID, days int) ([]TrendPoint, error) {
	if days <= 0 {
		days = 14
	}
	now := time.Now().UTC()
	start := now.AddDate(0, 0, -days+1)
	buckets := make(map[string]int64, days)
	for i := 0; i < days; i++ {
		d := start.AddDate(0, 0, i).Format("2006-01-02")
		buckets[d] = 0
	}
	rows, _, err := s.analyses.List(ctx, repositories.ListFilter{UserID: userID, Limit: 1000})
	if err != nil {
		return nil, err
	}
	for _, a := range rows {
		d := a.CreatedAt.UTC().Format("2006-01-02")
		if _, ok := buckets[d]; ok {
			buckets[d]++
		}
	}
	out := make([]TrendPoint, 0, days)
	for i := 0; i < days; i++ {
		d := start.AddDate(0, 0, i).Format("2006-01-02")
		out = append(out, TrendPoint{Date: d, Count: buckets[d]})
	}
	return out, nil
}

type RiskBucket struct {
	Level string `json:"level"`
	Count int64  `json:"count"`
}

func (s *DashboardService) RiskDistribution(ctx context.Context, userID uuid.UUID) ([]RiskBucket, error) {
	levels := []string{"low", "medium", "high", "critical"}
	out := make([]RiskBucket, 0, len(levels))
	for _, l := range levels {
		c, err := s.analyses.CountByRisk(ctx, userID, l)
		if err != nil {
			return nil, err
		}
		out = append(out, RiskBucket{Level: l, Count: c})
	}
	return out, nil
}

type TypeBucket struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

func (s *DashboardService) ObjectTypes(ctx context.Context, userID uuid.UUID) ([]TypeBucket, error) {
	types := []models.ObjectType{
		models.ObjectProcedure, models.ObjectFunction, models.ObjectPackage,
		models.ObjectTrigger, models.ObjectView, models.ObjectScript,
	}
	out := make([]TypeBucket, 0, len(types))
	for _, t := range types {
		c, err := s.analyses.CountByType(ctx, userID, string(t))
		if err != nil {
			return nil, err
		}
		out = append(out, TypeBucket{Type: string(t), Count: c})
	}
	return out, nil
}
