package endpoint

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func getCurrentUser(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return -1, errors.New("user_id not found from token")

	}

	userIDParsed, ok := userID.(int)
	if !ok {
		return -1, errors.New("cannot parse user_id")
	}

	return userIDParsed, nil
}
