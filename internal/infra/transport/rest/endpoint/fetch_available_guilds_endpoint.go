package endpoint

import (
	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/gin-gonic/gin"
)

type FetchAvailableGuildsResponse struct {
	Message string         `json:"message"`
	Data    []entity.Guild `json:"data"`
	Code    int            `json:"status_code"`
}

// Fetch available guilds godoc
//
//	@Summary		Fetch available guilds
//	@Description	Fetch available guilds
//	@Tags			Guilds
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	FetchAvailableGuildsResponse
//	@Failure		400		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/guilds [get]
func (h *Handler) fetchAvailableGuildsEndpoint(c *gin.Context) {
	page, err := connhelper.GetPage(c)
	if err != nil {
		connhelper.NewBadRequestError(c, err)
		return
	}

	params := core.PaginationParams{
		Page: page,
	}

	guilds, err := h.guildService.FetchAvailableGuilds(params)
	if err != nil {
		connhelper.NewInternalServerError(c, err)
		return
	}

	connhelper.NewOkResponse(c, "available guilds fetched", guilds)
}
