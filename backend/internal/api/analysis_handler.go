package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/YASSERRMD/sql-sage/backend/internal/middleware"
	"github.com/YASSERRMD/sql-sage/backend/internal/models"
	"github.com/YASSERRMD/sql-sage/backend/internal/repositories"
	"github.com/YASSERRMD/sql-sage/backend/internal/services"
	"github.com/YASSERRMD/sql-sage/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalysisHandler struct {
	svc *services.AnalysisService
}

func NewAnalysisHandler(s *services.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{svc: s}
}

func (h *AnalysisHandler) Register(rg *gin.RouterGroup, authMW gin.HandlerFunc) {
	g := rg.Group("/analyses", authMW, middleware.BodySizeLimit(int64(middleware.MaxBodyBytes)*4))
	g.POST("", h.create)
	g.GET("", h.list)
	g.GET("/:id", h.get)
	g.DELETE("/:id", h.delete)
}

func (h *AnalysisHandler) create(c *gin.Context) {
	var in services.CreateAnalysisInput
	if err := c.ShouldBindJSON(&in); err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", err.Error(), nil)
		return
	}
	uid := mustUserID(c)
	a, err := h.svc.Create(c.Request.Context(), uid, in)
	if err != nil {
		h.fail(c, err)
		return
	}
	out := map[string]any{
		"id":         a.ID,
		"providerId": a.ProviderID,
		"objectName": a.ObjectName,
		"objectType": a.ObjectType,
		"summary":    a.Summary,
		"riskScore":  a.RiskScore,
		"createdAt":  a.CreatedAt,
		"tokensUsed": a.TokensUsed,
	}
	schema, _ := h.svc.DecodeResult(a)
	if schema != nil {
		out["result"] = schema
	}
	httpx.Created(c, out)
}

func (h *AnalysisHandler) get(c *gin.Context) {
	uid := mustUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", "invalid id", nil)
		return
	}
	a, err := h.svc.Get(c.Request.Context(), uid, id)
	if err != nil {
		h.fail(c, err)
		return
	}
	out := map[string]any{
		"id":         a.ID,
		"providerId": a.ProviderID,
		"objectName": a.ObjectName,
		"objectType": a.ObjectType,
		"summary":    a.Summary,
		"riskScore":  a.RiskScore,
		"sourceCode": a.SourceCode,
		"createdAt":  a.CreatedAt,
		"tokensUsed": a.TokensUsed,
	}
	schema, _ := h.svc.DecodeResult(a)
	if schema != nil {
		out["result"] = schema
	}
	httpx.OK(c, out)
}

func (h *AnalysisHandler) delete(c *gin.Context) {
	uid := mustUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", "invalid id", nil)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uid, id); err != nil {
		h.fail(c, err)
		return
	}
	httpx.NoContent(c)
}

func (h *AnalysisHandler) list(c *gin.Context) {
	uid := mustUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	f := repositories.ListFilter{
		UserID:     uid,
		Query:      c.Query("q"),
		ObjectType: c.Query("objectType"),
		Risk:       c.Query("risk"),
		Limit:      size,
		Offset:     (page - 1) * size,
	}
	items, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		h.fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"items": items, "total": total, "page": page, "pageSize": size})
}

func (h *AnalysisHandler) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrAnalysisNotFound):
		httpx.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "analysis not found", nil)
	case errors.Is(err, services.ErrValidation):
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", err.Error(), nil)
	default:
		httpx.AbortWithError(c, http.StatusBadGateway, "PROVIDER", err.Error(), nil)
	}
}

var _ = json.Marshal
var _ = models.ObjectProcedure
