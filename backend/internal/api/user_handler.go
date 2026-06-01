package api

import (
	"net/http"
	"strings"

	"github.com/YASSERRMD/sql-sage/backend/internal/middleware"
	"github.com/YASSERRMD/sql-sage/backend/internal/models"
	"github.com/YASSERRMD/sql-sage/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler { return &UserHandler{db: db} }

type registerReq struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", err.Error(), nil)
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	var count int64
	if err := h.db.WithContext(c.Request.Context()).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		httpx.AbortWithError(c, http.StatusInternalServerError, "INTERNAL", err.Error(), nil)
		return
	}
	if count > 0 {
		httpx.AbortWithError(c, http.StatusConflict, "CONFLICT", "email already registered", nil)
		return
	}
	hash, err := hashPassword(req.Password)
	if err != nil {
		httpx.AbortWithError(c, http.StatusInternalServerError, "INTERNAL", err.Error(), nil)
		return
	}
	u := &models.User{
		Email:        email,
		Name:         strings.TrimSpace(req.Name),
		PasswordHash: hash,
		Role:         models.RoleUser,
		IsActive:     true,
	}
	if err := h.db.WithContext(c.Request.Context()).Create(u).Error; err != nil {
		httpx.AbortWithError(c, http.StatusInternalServerError, "INTERNAL", err.Error(), nil)
		return
	}
	httpx.Created(c, u)
}

type updateMeReq struct {
	Name string `json:"name"`
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	id, _ := uuid.Parse(middleware.GetUserID(c))
	var u models.User
	if err := h.db.WithContext(c.Request.Context()).Where("id = ?", id).First(&u).Error; err != nil {
		httpx.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "user not found", nil)
		return
	}
	var req updateMeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.AbortWithError(c, http.StatusBadRequest, "VALIDATION", err.Error(), nil)
		return
	}
	if req.Name != "" {
		u.Name = strings.TrimSpace(req.Name)
	}
	if err := h.db.WithContext(c.Request.Context()).Save(&u).Error; err != nil {
		httpx.AbortWithError(c, http.StatusInternalServerError, "INTERNAL", err.Error(), nil)
		return
	}
	httpx.OK(c, u)
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup, authMW gin.HandlerFunc) {
	rg.POST("/auth/register", h.Register)
	g := rg.Group("/users", authMW)
	g.PATCH("/me", h.UpdateMe)
}
