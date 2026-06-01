package middleware

import (
	"net/http"
	"strings"

	"github.com/YASSERRMD/sql-sage/backend/internal/auth"
	"github.com/YASSERRMD/sql-sage/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

const (
	CtxUserID = "ss_user_id"
	CtxRole   = "ss_role"
	CtxEmail  = "ss_email"
)

func AuthRequired(jwtSvc *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			httpx.AbortWithError(c, http.StatusUnauthorized, "AUTH_MISSING", "missing authorization header", nil)
			return
		}
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			httpx.AbortWithError(c, http.StatusUnauthorized, "AUTH_FORMAT", "invalid authorization format", nil)
			return
		}
		claims, err := jwtSvc.ParseAccess(parts[1])
		if err != nil {
			httpx.AbortWithError(c, http.StatusUnauthorized, "AUTH_INVALID", "invalid or expired token", nil)
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Set(CtxEmail, claims.Email)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role, _ := c.Get(CtxRole)
		r, _ := role.(string)
		if _, ok := allowed[r]; !ok {
			httpx.AbortWithError(c, http.StatusForbidden, "FORBIDDEN", "insufficient role", nil)
			return
		}
		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	v, _ := c.Get(CtxUserID)
	s, _ := v.(string)
	return s
}
