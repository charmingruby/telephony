package endpoint

import (
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
	userID, err := getCurrentUser(c)
	if err != nil {
		NewInternalServerError(c, err)
		return
	}

	profile, err := h.userService.GetProfileByID(userID)
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
