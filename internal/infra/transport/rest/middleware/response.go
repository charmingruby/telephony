package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewResponse(c *gin.Context, code int, data any, message string) {
	res := Response{
		Message: message,
		Data:    data,
		Code:    code,
	}
	c.JSON(code, res)
}

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Code    int    `json:"status_code"`
}

func NewUnauthorized(c *gin.Context, msg string, data any) {
	NewResponse(c, http.StatusUnauthorized, data, msg)
}
