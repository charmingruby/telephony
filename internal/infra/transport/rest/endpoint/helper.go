package endpoint

import (
	"errors"
	"fmt"
	"strconv"

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

func getParamID(c *gin.Context, identifier string) (int, error) {
	paramID := c.Param(identifier)

	id, err := strconv.Atoi(paramID)
	if err != nil {
		return -1, fmt.Errorf("cannot parse `%s`", identifier)
	}

	return id, nil
}

func getPage(c *gin.Context) (int, error) {
	var page int

	pageParams := c.DefaultQuery("page", "1")
	if pageParams == "" {
		page = 1
	}

	convPage, err := strconv.Atoi(pageParams)
	if err != nil {
		return -1, fmt.Errorf("`%s` is not a valid page", pageParams)
	}

	if convPage <= 0 {
		page = 1
	} else {
		page = convPage
	}

	return page, nil
}
