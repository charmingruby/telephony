package endpoint

import (
	"github.com/charmingruby/telephony/internal/domain/user/dto"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

type CreateProfileRequest struct {
	DisplayName string `json:"display_name" binding:"required"`
	Bio         string `json:"bio" binding:"required"`
	UserID      int    `json:"user_id" binding:"required"`
}

// Create profile godoc
//
//	@Summary		Create an user profile
//	@Description	Create an user profile
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
	var req CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		NewPayloadError(c, err)
		return
	}

	dto := dto.CreateProfileDTO{
		DisplayName: req.DisplayName,
		Bio:         req.Bio,
		UserID:      req.UserID,
	}

	if err := h.userService.CreateProfile(dto); err != nil {
		notFoundErr, ok := err.(*validation.ErrNotFound)
		if ok {
			NewResourceNotFoundError(c, notFoundErr)
			return
		}

		conflictErr, ok := err.(*validation.ErrConflict)
		if ok {
			NewConflicError(c, conflictErr)
			return
		}

		validationErr, ok := err.(*validation.ErrValidation)
		if ok {
			NewEntityError(c, validationErr)
			return
		}

		NewInternalServerError(c, err)
		return
	}

	NewCreatedResponse(c, "user profile")
}
