package endpoint

import (
	"github.com/charmingruby/telephony/internal/domain/user/dto"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// Register godoc
//
//	@Summary		Create user
//	@Description	Create a new user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest	true	"Add User"
//	@Success		201		{object}	Response
//	@Failure		400		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/auth/register [post]
func (h *Handler) registerEndpoint(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newPayloadError(c, err)
		return
	}

	dto := dto.RegisterDTO{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}

	if err := h.userService.Register(dto); err != nil {
		validationErr, ok := err.(*validation.ErrValidation)
		if ok {
			newEntityError(c, validationErr)
			return
		}

		conflictErr, ok := err.(*validation.ErrConflict)
		if ok {
			newConflicError(c, conflictErr)
			return
		}

		newInternalServerError(c, err)
		return
	}
	newCreatedResponse(c, "user")
}
