package middleware

import (
	"strings"

	"github.com/charmingruby/telephony/internal/infra/security/token"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwt *token.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization") // Bearer <token>
		if bearerToken == "" {
			NewUnauthorized(c, "token not found on header: Authorization", nil)
			c.Abort()
			return
		}

		token := strings.Split(bearerToken, " ")[1]

		payload, err := jwt.ValidateToken(token)
		if err != nil {
			NewUnauthorized(c, err.Error(), nil)
			c.Abort()
			return
		}

		c.Set("user_id", payload.UserID)
		c.Next()
	}
}
