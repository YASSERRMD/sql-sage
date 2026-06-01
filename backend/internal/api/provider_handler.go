package api

import (
	"errors"
	"net/http"

	"github.com/YASSERRMD/sql-sage/backend/internal/middleware"
	"github.com/YASSERRMD/sql-sage/backend/internal/models"
	"github.com/YASSERRMD/sql-sage/backend/internal/services"
	"github.com/YASSERRMD/sql-sage/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProviderHandler struct {
	svc *services.ProviderService
}

func NewProviderHandler(s *services.ProviderService) *ProviderHandler {
	return &ProviderHandler{svc: s}
}

func (h *ProviderHandler) Register(rg *gin.RouterGroup, authMW gin.HandlerFunc) {
	g := rg.Group("/providers", authMW)
	g.GET("", h.list)
	g.POST("", h.create)
	g.GET("/:id", h.get)
	g.PUT("/:id", h.update)
	g.DELETE("/:id", h.delete)
	g.POST("/:id/test", h.test)
	g.POST("/:id/default", h.setDefault)
}

func (h *ProviderHandler) list(c *gin.Context) {
	uid := mustUserID(c)
	items, err := h.svc.List(c.Request.Context(), uid)
	if err != nil {
		httpx.AbortWithError(c, http.StatusInternalServerError, "INTERNAL", err.Error(), nil)
		return
	}
	httpx.OK(c, items)
}

func (h *ProviderHandler) create(c *gin.Context) {
	var in services.CreateProviderInput
	if err := c.ShouldBindJSON(&in); err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", err.Error(), nil)
		return
	}
	uid := mustUserID(c)
	p, err := h.svc.Create(c.Request.Context(), uid, in)
	if err != nil {
		h.fail(c, err)
		return
	}
	httpx.Created(c, p)
}

func (h *ProviderHandler) get(c *gin.Context) {
	uid := mustUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", "invalid id", nil)
		return
	}
	p, err := h.svc.Get(c.Request.Context(), uid, id)
	if err != nil {
		h.fail(c, err)
		return
	}
	httpx.OK(c, p)
}

func (h *ProviderHandler) update(c *gin.Context) {
	uid := mustUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", "invalid id", nil)
		return
	}
	var in services.CreateProviderInput
	if err := c.ShouldBindJSON(&in); err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", err.Error(), nil)
		return
	}
	p, err := h.svc.Update(c.Request.Context(), uid, id, in)
	if err != nil {
		h.fail(c, err)
		return
	}
	httpx.OK(c, p)
}

func (h *ProviderHandler) delete(c *gin.Context) {
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

func (h *ProviderHandler) test(c *gin.Context) {
	uid := mustUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", "invalid id", nil)
		return
	}
	res, err := h.svc.Test(c.Request.Context(), uid, id)
	if err != nil {
		httpx.AbortWithError(c, http.StatusBadGateway, "PROVIDER", err.Error(), nil)
		return
	}
	httpx.OK(c, res)
}

func (h *ProviderHandler) setDefault(c *gin.Context) {
	uid := mustUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", "invalid id", nil)
		return
	}
	p, err := h.svc.SetDefault(c.Request.Context(), uid, id)
	if err != nil {
		h.fail(c, err)
		return
	}
	httpx.OK(c, p)
}

func (h *ProviderHandler) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrProviderNotFound):
		httpx.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "provider not found", nil)
	case errors.Is(err, services.ErrValidation):
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", err.Error(), nil)
	default:
		httpx.AbortWithError(c, http.StatusInternalServerError, "INTERNAL", err.Error(), nil)
	}
}

func mustUserID(c *gin.Context) uuid.UUID {
	id, _ := uuid.Parse(middleware.GetUserID(c))
	return id
}

var _ = models.RoleAdmin
