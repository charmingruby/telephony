package endpoint

import (
	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

type FetchGuildChannelsRequest struct {
	ProfileID int `json:"profile_id"`
}

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
func (h *Handler) fetchGuildChannelsEndpoint(c *gin.Context) {
	guildID, err := getParamID(c, "guild_id")
	if err != nil {
		NewBadRequestError(c, err)
		return
	}

	userID, err := getCurrentUser(c)
	if err != nil {
		NewInternalServerError(c, err)
		return
	}

	page, err := getPage(c)
	if err != nil {
		NewBadRequestError(c, err)
		return
	}

	var req FetchGuildChannelsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		NewPayloadError(c, err)
		return
	}

	dto := dto.FetchGuildChannelsDTO{
		UserID:    userID,
		ProfileID: req.ProfileID,
		GuildID:   guildID,
		Pagination: core.PaginationParams{
			Page: page,
		},
	}

	channels, err := h.guildService.FetchGuildChannels(dto)
	if err != nil {
		notFoundErr, ok := err.(*validation.ErrNotFound)
		if ok {
			NewResourceNotFoundError(c, notFoundErr)
			return
		}

		unauthorizedErr, ok := err.(*validation.ErrUnathorized)
		if ok {
			NewUnauthorizedErr(c, unauthorizedErr.Error())
			return
		}

		NewInternalServerError(c, err)
		return
	}

	NewOkResponse(c, "guild channels fetched", channels)
}
