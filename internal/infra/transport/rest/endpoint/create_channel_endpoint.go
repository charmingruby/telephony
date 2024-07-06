package endpoint

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

type CreateChannelRequest struct {
	Name      string `json:"name" binding:"required"`
	ProfileID int    `json:"profile_id" binding:"required"`
}

// Create guild godoc
//
//	@Summary		Creates a channel
//	@Description	Creates a channel
//	@Tags			Channels
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateChannelRequest	true	"Create Channel Payload"
//	@Success		201		{object}	Response
//	@Failure		401		{object}	Response
//	@Failure		404		{object}	Response
//	@Failure		409		{object}	Response
//	@Failure		422		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/guilds/{guild_id}/channels [post]
func (h *Handler) createChannelEndpoint(c *gin.Context) {
	userID, err := connhelper.GetCurrentUser(c)
	if err != nil {
		connhelper.NewInternalServerError(c, err)
		return
	}

	guildID, err := connhelper.GetParamID(c, "guild_id")
	if err != nil {
		connhelper.NewInternalServerError(c, err)
		return
	}

	var req CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		connhelper.NewPayloadError(c, err)
		return
	}

	dto := dto.CreateChannelDTO{
		Name:      req.Name,
		GuildID:   guildID,
		ProfileID: req.ProfileID,
		UserID:    userID,
	}

	if _, err := h.guildService.CreateChannel(dto); err != nil {
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

	connhelper.NewCreatedResponse(c, "channel")
}
