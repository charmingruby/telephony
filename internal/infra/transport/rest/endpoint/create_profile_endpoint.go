package endpoint

import (
	"github.com/charmingruby/telephony/internal/domain/user/dto"
	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

type CreateProfileRequest struct {
	DisplayName string `json:"display_name" binding:"required"`
	Bio         string `json:"bio" binding:"required"`
}

// Create profile godoc
//
//	@Summary		Creates an user profile
//	@Description	Creates an user profile
//	@Tags			Profiles
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateProfileRequest	true	"Create Profile Payload"
//	@Success		201		{object}	Response
//	@Failure		422		{object}	Response
//	@Failure		404		{object}	Response
//	@Failure		409		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/me/profile [post]
func (h *Handler) createProfileEndpoint(c *gin.Context) {
	userID, err := connhelper.GetCurrentUser(c)
	if err != nil {
		connhelper.NewInternalServerError(c, err)
		return
	}

	var req CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		connhelper.NewPayloadError(c, err)
		return
	}

	dto := dto.CreateProfileDTO{
		DisplayName: req.DisplayName,
		Bio:         req.Bio,
		UserID:      userID,
	}

	if err := h.userService.CreateProfile(dto); err != nil {
		notFoundErr, ok := err.(*validation.ErrNotFound)
		if ok {
			connhelper.NewResourceNotFoundError(c, notFoundErr)
			return
		}

		conflictErr, ok := err.(*validation.ErrConflict)
		if ok {
			connhelper.NewConflicError(c, conflictErr)
			return
		}

		validationErr, ok := err.(*validation.ErrValidation)
		if ok {
			connhelper.NewEntityError(c, validationErr)
			return
		}

		connhelper.NewInternalServerError(c, err)
		return
	}

	connhelper.NewCreatedResponse(c, "user profile")
}
