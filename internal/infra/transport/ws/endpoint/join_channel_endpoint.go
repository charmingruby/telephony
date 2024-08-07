package endpoint

import (
	"strconv"

	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/charmingruby/telephony/internal/infra/transport/ws"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

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

	strProfileID, err := connhelper.GetQueryParamValue(c, "profile")
	if err != nil {
		connhelper.NewBadRequestError(c, err)
		return
	}

	profileID, err := strconv.Atoi(strProfileID)
	if err != nil {
		connhelper.NewBadRequestError(c, err)
		return
	}

	dto := dto.JoinChannelDTO{
		ChannelID: channelID,
		GuildID:   guildID,
		ProfileID: profileID,
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
