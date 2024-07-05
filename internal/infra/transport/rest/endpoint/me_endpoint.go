package endpoint

import (
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
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
//	@Summary		Gets authenticated user profile
//	@Description	Gets authenticated user profile
//	@Tags			Profiles
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	MeResponse
//	@Failure		404		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/me [get]
func (h *Handler) meEndpoint(c *gin.Context) {
	userID, err := connhelper.GetCurrentUser(c)
	if err != nil {
		connhelper.NewInternalServerError(c, err)
		return
	}

	profile, err := h.userService.GetProfileByID(userID)
	if err != nil {
		notFoundErr, ok := err.(*validation.ErrNotFound)
		if ok {
			connhelper.NewResourceNotFoundError(c, notFoundErr)
			return
		}

		connhelper.NewInternalServerError(c, err)
		return
	}

	connhelper.NewOkResponse(
		c,
		"profile found",
		profile,
	)
}
