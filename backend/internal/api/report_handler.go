package api

import (
	"encoding/json"
	"net/http"

	"github.com/YASSERRMD/sql-sage/backend/internal/analysis"
	"github.com/YASSERRMD/sql-sage/backend/internal/middleware"
	"github.com/YASSERRMD/sql-sage/backend/internal/report"
	"github.com/YASSERRMD/sql-sage/backend/internal/services"
	"github.com/YASSERRMD/sql-sage/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *AnalysisHandler) Report(c *gin.Context) {
	uid := mustUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", "invalid id", nil)
		return
	}
	a, err := h.svc.Get(c.Request.Context(), uid, id)
	if err != nil {
		httpx.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "analysis not found", nil)
		return
	}
	schema, err := h.svc.DecodeResult(a)
	if err != nil {
		httpx.AbortWithError(c, http.StatusInternalServerError, "INTERNAL", "invalid stored result", nil)
		return
	}
	in := report.Input{
		ObjectName:   a.ObjectName,
		ObjectType:   string(a.ObjectType),
		Summary:      a.Summary,
		MarkdownBody: schema.MarkdownReport,
		Mermaid:      schema.MermaidDiagram,
		CreatedAt:    a.CreatedAt,
		Risk:         a.RiskScore,
		TokensUsed:   a.TokensUsed,
	}
	format := c.DefaultQuery("format", "md")
	switch format {
	case "md", "markdown":
		c.Header("Content-Type", "text/markdown; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=\""+a.ObjectName+".md\"")
		c.String(http.StatusOK, report.Markdown(in))
	case "html":
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=\""+a.ObjectName+".html\"")
		c.String(http.StatusOK, report.HTML(in))
	case "pdf":
		httpx.AbortWithError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "PDF export not yet implemented", nil)
	default:
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", "format must be md, html, or pdf", nil)
	}
}

var _ = analysis.SystemPrompt
var _ = json.Marshal
var _ = services.ErrAnalysisNotFound
var _ = middleware.CtxUserID
