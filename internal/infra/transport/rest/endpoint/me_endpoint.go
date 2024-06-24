package endpoint

import (
	"fmt"

	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

type MeResponse struct {
	Message string              `json:"message"`
	Data    *entity.UserProfile `json:"data"`
	Code    int                 `json:"status_code"`
}

// Me godoc
//
//	@Summary		Get authenticated user profile
//	@Description	Get authenticated user profile
//	@Tags			Profiles
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	MeResponse
//	@Failure		404		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/me [get]
func (h *Handler) meEndpoint(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		NewUnauthorizedErr(c, "user_id not found from token")
		return
	}

	userIDParsed, ok := userID.(int)
	if !ok {
		NewInternalServerError(c, fmt.Errorf("cannot parse user_id"))
		return
	}

	profile, err := h.userService.GetProfileByID(userIDParsed)
	if err != nil {
		notFoundErr, ok := err.(*validation.ErrNotFound)
		if ok {
			NewResourceNotFoundError(c, notFoundErr)
			return
		}

		NewInternalServerError(c, err)
		return
	}

	NewOkResponse(
		c,
		"profile found",
		profile,
	)
}
