package analysis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/YASSERRMD/sql-sage/backend/internal/models"
	"github.com/YASSERRMD/sql-sage/backend/pkg/crypto"
	"github.com/YASSERRMD/sql-sage/backend/pkg/llm"
	"github.com/google/uuid"
)

const maxRetries = 2

type ProviderResolver interface {
	Get(ctx context.Context, userID, id uuid.UUID) (*models.Provider, error)
	FindDefault(ctx context.Context, userID uuid.UUID) (*models.Provider, error)
}

type AnalysisStore interface {
	Create(ctx context.Context, a *models.Analysis) error
}

type Engine struct {
	llm      *llm.Client
	cipher   *crypto.Cipher
	resolver ProviderResolver
	store    AnalysisStore
}

func NewEngine(l *llm.Client, c *crypto.Cipher, r ProviderResolver, s AnalysisStore) *Engine {
	return &Engine{llm: l, cipher: c, resolver: r, store: s}
}

type AnalyzeInput struct {
	UserID     uuid.UUID
	ProviderID *uuid.UUID
	ObjectName string
	ObjectType models.ObjectType
	SourceCode string
}

type AnalyzeResult struct {
	AnalysisID uuid.UUID
	Schema     *Schema
	TokensUsed int
}

func (e *Engine) Analyze(ctx context.Context, in AnalyzeInput) (*AnalyzeResult, error) {
	p, err := e.resolveProvider(ctx, in.UserID, in.ProviderID)
	if err != nil {
		return nil, err
	}
	apiKey, err := e.cipher.Decrypt(p.APIKeyEncrypted)
	if err != nil {
		return nil, err
	}

	messages := []llm.ChatMessage{
		{Role: "system", Content: SystemPrompt},
		{Role: "user", Content: BuildUserPrompt(in.ObjectName, string(in.ObjectType), in.SourceCode)},
	}

	schema, used, err := e.callWithRetry(ctx, p, apiKey, messages)
	if err != nil {
		return nil, err
	}
	a := &models.Analysis{
		UserID:     in.UserID,
		ProviderID: p.ID,
		ObjectName: firstNonEmpty(schema.ObjectName, in.ObjectName),
		ObjectType: models.ObjectType(firstNonEmpty(schema.ObjectType, string(in.ObjectType))),
		SourceCode: in.SourceCode,
		Summary:    schema.Summary,
		RiskScore:  deriveRisk(schema),
		Result:     marshalSchema(schema),
		TokensUsed: used,
	}
	if !a.ObjectType.Valid() {
		a.ObjectType = models.ObjectUnknown
	}
	if err := e.store.Create(ctx, a); err != nil {
		return nil, fmt.Errorf("store analysis: %w", err)
	}
	return &AnalyzeResult{AnalysisID: a.ID, Schema: schema, TokensUsed: used}, nil
}

func (e *Engine) callWithRetry(
	ctx context.Context,
	p *models.Provider,
	apiKey string,
	messages []llm.ChatMessage,
) (*Schema, int, error) {
	total := 0
	current := messages
	for attempt := 0; attempt <= maxRetries; attempt++ {
		cctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
		resp, err := e.llm.Chat(cctx, llm.ChatRequest{
			BaseURL:     p.BaseURL,
			APIKey:      apiKey,
			Model:       p.ModelName,
			Temperature: p.Temperature,
			MaxTokens:   p.MaxTokens,
			Messages:    current,
		})
		cancel()
		if err != nil {
			if attempt == maxRetries {
				return nil, total, err
			}
			continue
		}
		total += resp.Usage
		schema, verr := ValidateAndParse(resp.Content)
		if verr == nil {
			return schema, total, nil
		}
		if attempt == maxRetries {
			return nil, total, verr
		}
		repaired := RepairJSON(resp.Content)
		current = []llm.ChatMessage{
			current[0],
			{Role: "user", Content: BuildRepairUserPrompt(repaired, verr.Error())},
		}
	}
	return nil, total, errors.New("analysis failed")
}

func (e *Engine) resolveProvider(ctx context.Context, userID uuid.UUID, id *uuid.UUID) (*models.Provider, error) {
	if id != nil {
		return e.resolver.Get(ctx, userID, *id)
	}
	return e.resolver.FindDefault(ctx, userID)
}

func marshalSchema(s *Schema) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func firstNonEmpty(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

func deriveRisk(s *Schema) string {
	highest := "low"
	rank := map[string]int{"low": 1, "medium": 2, "high": 3, "critical": 4}
	for _, r := range s.Risks {
		m, _ := r.(map[string]any)
		lvl, _ := m["level"].(string)
		if rank[lvl] > rank[highest] {
			highest = lvl
		}
	}
	return highest
}
