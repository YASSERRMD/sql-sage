package httpx

import (
	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

func AbortWithError(c *gin.Context, status int, code, message string, details map[string]any) {
	c.AbortWithStatusJSON(status, ErrorResponse{
		Error: ErrorBody{Code: code, Message: message, Details: details},
	})
}

func OK(c *gin.Context, data any) {
	c.JSON(200, data)
}

func Created(c *gin.Context, data any) {
	c.JSON(201, data)
}

func NoContent(c *gin.Context) {
	c.Status(204)
}
