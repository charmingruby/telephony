package endpoint

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

type CreateGuildRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	ProfileID   int    `json:"profile_id" binding:"required"`
}

// Create guild godoc
//
//	@Summary		Creates a guild
//	@Description	Creates a guild
//	@Tags			Guilds
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateGuildRequest	true	"Create Guild Payload"
//	@Success		201		{object}	Response
//	@Failure		401		{object}	Response
//	@Failure		404		{object}	Response
//	@Failure		409		{object}	Response
//	@Failure		422		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/guilds [post]
func (h *Handler) createGuildEndpoint(c *gin.Context) {
	userID, err := connhelper.GetCurrentUser(c)
	if err != nil {
		connhelper.NewInternalServerError(c, err)
		return
	}

	var req CreateGuildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		connhelper.NewPayloadError(c, err)
		return
	}

	dto := dto.CreateGuildDTO{
		Name:        req.Name,
		Description: req.Description,
		ProfileID:   req.ProfileID,
		UserID:      userID,
	}

	if err := h.guildService.CreateGuild(dto); err != nil {
		notFoundErr, ok := err.(*validation.ErrNotFound)
		if ok {
			connhelper.NewResourceNotFoundError(c, notFoundErr)
			return
		}

		unauthorizedErr, ok := err.(*validation.ErrUnathorized)
		if ok {
			connhelper.NewUnauthorizedErr(c, unauthorizedErr.Error())
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

	connhelper.NewCreatedResponse(c, "guild")
}
