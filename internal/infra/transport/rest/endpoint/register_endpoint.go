package endpoint

import (
	"github.com/charmingruby/telephony/internal/domain/user/dto"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name string `json:"name" binding:"required"`
}

// Register godoc
//
//	@Summary		Create example
//	@Description	Create a new example
//	@Tags			Examples
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest	true	"Add Example"
//	@Success		201		{object}	Response
//	@Failure		400		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/examples [post]
func (h *Handler) RegisterEndpoint(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newPayloadError(c, err)
		return
	}

	dto := dto.RegisterDTO{}

	if err := h.userService.Register(dto); err != nil {
		validationErr, ok := err.(*validation.ErrValidation)
		if ok {
			newBadRequestError(c, validationErr)
			return
		}

		newInternalServerError(c, err)
		return
	}
	newCreatedResponse(c, "example")
}
