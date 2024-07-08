package endpoint

import (
	"fmt"
	"strconv"

	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/charmingruby/telephony/internal/infra/transport/ws/presenter"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

func (h *WebSocketHandler) getConnectedClientsEndpoint(c *gin.Context) {
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

	channelID, err := connhelper.GetParamID(c, "channel_id")
	if err != nil {
		connhelper.NewInternalServerError(c, err)
		return
	}

	dto := dto.ValidatePermissionToCheckConnectedProfilesOnChannelDTO{
		ProfileID: profileID,
		UserID:    userID,
		GuildID:   guildID,
	}

	if err := h.guildService.ValidatePermissionToCheckConnectedProfilesOnChannel(dto); err != nil {
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

	var clients []presenter.ClientHTTP

	_, ok := h.hub.Rooms[channelID]
	if !ok {
		connhelper.NewResourceNotFoundError(c, fmt.Errorf("`%d` channel not found", channelID))
		return
	}

	for _, c := range h.hub.Rooms[channelID].Clients {
		clients = append(clients, presenter.ToHTTP(*c))
	}

	connhelper.NewOkResponse(c, "clients fetched", clients)
}
