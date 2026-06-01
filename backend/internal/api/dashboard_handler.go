package api

import (
	"net/http"

	"github.com/YASSERRMD/sql-sage/backend/internal/middleware"
	"github.com/YASSERRMD/sql-sage/backend/internal/services"
	"github.com/YASSERRMD/sql-sage/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	svc *services.DashboardService
}

func NewDashboardHandler(s *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: s}
}

func (h *DashboardHandler) Register(rg *gin.RouterGroup, authMW gin.HandlerFunc) {
	g := rg.Group("/dashboard", authMW)
	g.GET("/summary", h.summary)
	g.GET("/trend", h.trend)
	g.GET("/risk-distribution", h.risk)
	g.GET("/object-types", h.types)
}

func (h *DashboardHandler) summary(c *gin.Context) {
	uid, _ := uuid.Parse(middleware.GetUserID(c))
	s, err := h.svc.Summary(c.Request.Context(), uid)
	if err != nil {
		httpx.AbortWithError(c, http.StatusInternalServerError, "INTERNAL", err.Error(), nil)
		return
	}
	httpx.OK(c, s)
}

func (h *DashboardHandler) trend(c *gin.Context) {
	uid, _ := uuid.Parse(middleware.GetUserID(c))
	t, err := h.svc.Trend(c.Request.Context(), uid, 14)
	if err != nil {
		httpx.AbortWithError(c, http.StatusInternalServerError, "INTERNAL", err.Error(), nil)
		return
	}
	httpx.OK(c, t)
}

func (h *DashboardHandler) risk(c *gin.Context) {
	uid, _ := uuid.Parse(middleware.GetUserID(c))
	r, err := h.svc.RiskDistribution(c.Request.Context(), uid)
	if err != nil {
		httpx.AbortWithError(c, http.StatusInternalServerError, "INTERNAL", err.Error(), nil)
		return
	}
	httpx.OK(c, r)
}

func (h *DashboardHandler) types(c *gin.Context) {
	uid, _ := uuid.Parse(middleware.GetUserID(c))
	t, err := h.svc.ObjectTypes(c.Request.Context(), uid)
	if err != nil {
		httpx.AbortWithError(c, http.StatusInternalServerError, "INTERNAL", err.Error(), nil)
		return
	}
	httpx.OK(c, t)
}
