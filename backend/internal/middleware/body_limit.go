package middleware

import (
	"net/http"

	"github.com/YASSERRMD/sql-sage/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

const MaxBodyBytes = 1 << 20 // 1 MiB

func BodySizeLimit(max int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > max {
			httpx.AbortWithError(c, http.StatusRequestEntityTooLarge, "TOO_LARGE", "request body too large", nil)
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, max)
		c.Next()
	}
}
