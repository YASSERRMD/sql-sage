package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/YASSERRMD/sql-sage/backend/internal/models"
	"github.com/YASSERRMD/sql-sage/backend/internal/repositories"
	"github.com/YASSERRMD/sql-sage/backend/pkg/crypto"
	"github.com/YASSERRMD/sql-sage/backend/pkg/llm"
	"github.com/google/uuid"
)

var (
	ErrProviderNotFound = errors.New("provider not found")
	ErrValidation       = errors.New("validation failed")
)

type ProviderService struct {
	repo      *repositories.ProviderRepository
	cipher    *crypto.Cipher
	llm       *llm.Client
	allowlist []string
}

func NewProviderService(r *repositories.ProviderRepository, c *crypto.Cipher, l *llm.Client, allowlist []string) *ProviderService {
	return &ProviderService{repo: r, cipher: c, llm: l, allowlist: allowlist}
}

type CreateProviderInput struct {
	Name        string  `json:"name"`
	BaseURL     string  `json:"baseUrl"`
	APIKey      string  `json:"apiKey"`
	ModelName   string  `json:"modelName"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"maxTokens"`
	IsDefault   bool    `json:"isDefault"`
}

func (s *ProviderService) Create(ctx context.Context, userID uuid.UUID, in CreateProviderInput) (*models.Provider, error) {
	if err := s.validate(in); err != nil {
		return nil, err
	}
	if err := s.llm.ValidateURL(in.BaseURL, s.allowlist); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidation, err.Error())
	}
	enc, err := s.cipher.Encrypt(in.APIKey)
	if err != nil {
		return nil, err
	}
	if in.IsDefault {
		_ = s.repo.ClearDefault(ctx, userID)
	}
	p := &models.Provider{
		UserID:          userID,
		Name:            in.Name,
		BaseURL:         in.BaseURL,
		APIKeyEncrypted: enc,
		APIKeyPreview:   mask(in.APIKey),
		ModelName:       in.ModelName,
		Temperature:     in.Temperature,
		MaxTokens:       in.MaxTokens,
		IsDefault:       in.IsDefault,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProviderService) Update(ctx context.Context, userID, id uuid.UUID, in CreateProviderInput) (*models.Provider, error) {
	p, err := s.repo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, ErrProviderNotFound
	}
	if in.Name != "" {
		p.Name = in.Name
	}
	if in.BaseURL != "" {
		if err := s.llm.ValidateURL(in.BaseURL, s.allowlist); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrValidation, err.Error())
		}
		p.BaseURL = in.BaseURL
	}
	if in.APIKey != "" {
		enc, err := s.cipher.Encrypt(in.APIKey)
		if err != nil {
			return nil, err
		}
		p.APIKeyEncrypted = enc
		p.APIKeyPreview = mask(in.APIKey)
	}
	if in.ModelName != "" {
		p.ModelName = in.ModelName
	}
	if in.Temperature > 0 {
		p.Temperature = in.Temperature
	}
	if in.MaxTokens > 0 {
		p.MaxTokens = in.MaxTokens
	}
	if in.IsDefault && !p.IsDefault {
		_ = s.repo.ClearDefault(ctx, userID)
	}
	p.IsDefault = in.IsDefault
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProviderService) Delete(ctx context.Context, userID, id uuid.UUID) error {
	return s.repo.Delete(ctx, userID, id)
}

func (s *ProviderService) Get(ctx context.Context, userID, id uuid.UUID) (*models.Provider, error) {
	p, err := s.repo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, ErrProviderNotFound
	}
	return p, nil
}

func (s *ProviderService) List(ctx context.Context, userID uuid.UUID) ([]models.Provider, error) {
	return s.repo.List(ctx, userID)
}

func (s *ProviderService) Test(ctx context.Context, userID, id uuid.UUID) (*llm.TestResult, error) {
	p, err := s.repo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, ErrProviderNotFound
	}
	key, err := s.cipher.Decrypt(p.APIKeyEncrypted)
	if err != nil {
		return nil, err
	}
	return s.llm.TestConnection(ctx, llm.TestRequest{
		BaseURL:       p.BaseURL,
		APIKey:        key,
		Model:         p.ModelName,
		HostAllowlist: s.allowlist,
	})
}

func (s *ProviderService) SetDefault(ctx context.Context, userID, id uuid.UUID) (*models.Provider, error) {
	if err := s.repo.SetDefault(ctx, userID, id); err != nil {
		return nil, ErrProviderNotFound
	}
	return s.repo.FindByID(ctx, userID, id)
}

func (s *ProviderService) validate(in CreateProviderInput) error {
	if strings.TrimSpace(in.Name) == "" {
		return fmt.Errorf("%w: name required", ErrValidation)
	}
	if strings.TrimSpace(in.BaseURL) == "" {
		return fmt.Errorf("%w: baseUrl required", ErrValidation)
	}
	if strings.TrimSpace(in.APIKey) == "" {
		return fmt.Errorf("%w: apiKey required", ErrValidation)
	}
	if strings.TrimSpace(in.ModelName) == "" {
		return fmt.Errorf("%w: modelName required", ErrValidation)
	}
	if in.Temperature < 0 || in.Temperature > 2 {
		return fmt.Errorf("%w: temperature must be 0..2", ErrValidation)
	}
	if in.MaxTokens <= 0 {
		return fmt.Errorf("%w: maxTokens must be > 0", ErrValidation)
	}
	return nil
}

func mask(k string) string {
	if len(k) <= 8 {
		return "****"
	}
	return k[:4] + "..." + k[len(k)-4:]
}
