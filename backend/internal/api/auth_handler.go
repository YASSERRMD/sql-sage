package api

import (
	"net/http"

	"github.com/YASSERRMD/sql-sage/backend/internal/middleware"
	"github.com/YASSERRMD/sql-sage/backend/internal/services"
	"github.com/YASSERRMD/sql-sage/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler { return &AuthHandler{svc: s} }

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type refreshReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type logoutReq struct {
	RefreshToken string `json:"refreshToken"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", err.Error(), nil)
		return
	}
	res, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials, services.ErrUserInactive:
			httpx.AbortWithError(c, http.StatusUnauthorized, "AUTH_INVALID", "invalid credentials", nil)
		default:
			httpx.AbortWithError(c, http.StatusInternalServerError, "INTERNAL", err.Error(), nil)
		}
		return
	}
	httpx.OK(c, gin.H{
		"accessToken":  res.AccessToken,
		"refreshToken": res.RefreshToken,
		"expiresIn":    res.ExpiresIn,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", err.Error(), nil)
		return
	}
	res, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		httpx.AbortWithError(c, http.StatusUnauthorized, "AUTH_INVALID", "invalid refresh token", nil)
		return
	}
	httpx.OK(c, gin.H{
		"accessToken":  res.AccessToken,
		"refreshToken": res.RefreshToken,
		"expiresIn":    res.ExpiresIn,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req logoutReq
	_ = c.ShouldBindJSON(&req)
	_ = h.svc.Logout(c.Request.Context(), req.RefreshToken)
	httpx.NoContent(c)
}

func (h *AuthHandler) Me(c *gin.Context) {
	idStr := middleware.GetUserID(c)
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpx.AbortWithError(c, http.StatusUnauthorized, "AUTH_INVALID", "invalid user id", nil)
		return
	}
	u, err := h.svc.Me(c.Request.Context(), id)
	if err != nil {
		httpx.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "user not found", nil)
		return
	}
	httpx.OK(c, u)
}

func (h *AuthHandler) Register(rg *gin.RouterGroup, authMW gin.HandlerFunc) {
	rg.POST("/auth/login", h.Login)
	rg.POST("/auth/refresh", h.Refresh)
	rg.POST("/auth/logout", h.Logout)
	rg.GET("/auth/me", authMW, h.Me)
}
