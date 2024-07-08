package endpoint

import (
	"strconv"

	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

type FetchGuildChannelsResponse struct {
	Message string           `json:"message"`
	Data    []entity.Channel `json:"data"`
	Code    int              `json:"status_code"`
}

// Fetch guild channels godoc
//
//	@Summary		Fetch paginated channels of a guild
//	@Description	Fetch paginated channels of a guild
//	@Tags			Channels
//	@Accept			json
//	@Produce		json
//	@Param			request	body		FetchGuildChannelsRequest	true	"Fetch Guild Channels Payload"
//
// @Success		200		{object}	FetchGuildChannelsResponse
// @Failure		400		{object}	Response
// @Failure		404		{object}	Response
// @Failure		500		{object}	Response
// @Router			/guilds/{guild_id}/channels [get]
func (h *WebSocketHandler) getAvailableGuildChannelsEndpoint(c *gin.Context) {
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

	dto := dto.FetchGuildChannelsDTO{
		UserID:    userID,
		ProfileID: profileID,
		GuildID:   guildID,
	}

	channels, err := h.guildService.FetchGuildChannels(dto)
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

	connhelper.NewOkResponse(c, "channels fetched", channels)
}
