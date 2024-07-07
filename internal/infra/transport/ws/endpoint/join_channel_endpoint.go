package endpoint

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/charmingruby/telephony/internal/infra/transport/ws"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

type JoinChannelRequest struct {
	ProfileID int `json:"profile_id" binding:"required"`
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
func (h *WebSocketHandler) joinChannelEndpoint(c *gin.Context) {
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

	channelID, err := connhelper.GetParamID(c, "channel_id")
	if err != nil {
		connhelper.NewInternalServerError(c, err)
		return
	}

	// var req JoinChannelRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	connhelper.NewPayloadError(c, err)
	// 	return
	// }

	dto := dto.JoinChannelDTO{
		ChannelID: channelID,
		GuildID:   guildID,
		ProfileID: 2,
		UserID:    userID,
	}

	res, err := h.guildService.JoinChannel(dto)

	if err != nil {
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

	client, err := ws.NewClient(c, h.hub, res.DisplayName, dto.ProfileID, dto.ChannelID, dto.GuildID)
	if err != nil {
		connhelper.NewInternalServerError(c, err)
		return
	}

	go client.WriteMessage()
	client.ReadMessage(h.hub)

	connhelper.NewOkResponse(c, "message sent", nil)
}
