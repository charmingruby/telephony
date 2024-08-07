package endpoint

import (
	"github.com/charmingruby/telephony/internal/domain/user/dto"
	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
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
//	@Summary		Creates an user
//	@Description	Creates an user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest	true	"Create User Payload"
//	@Success		201		{object}	Response
//	@Failure		422		{object}	Response
//	@Failure		409		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/auth/register [post]
func (h *Handler) registerEndpoint(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		connhelper.NewPayloadError(c, err)
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
			connhelper.NewEntityError(c, validationErr)
			return
		}

		conflictErr, ok := err.(*validation.ErrConflict)
		if ok {
			connhelper.NewConflicError(c, conflictErr)
			return
		}

		connhelper.NewInternalServerError(c, err)
		return
	}
	connhelper.NewCreatedResponse(c, "user")
}
