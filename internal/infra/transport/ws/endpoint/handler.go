package endpoint

import (
	"github.com/charmingruby/telephony/internal/domain/guild/usecase"
	"github.com/charmingruby/telephony/internal/infra/security/token"
	"github.com/charmingruby/telephony/internal/infra/transport/common/middleware"
	"github.com/charmingruby/telephony/internal/infra/transport/ws"
	"github.com/gin-gonic/gin"
)

func NewWSHandler(
	router *gin.Engine,
	guildService usecase.GuildServiceContract,
	token *token.JWTService,
	hub *ws.Hub,
) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

type WebSocketHandler struct {
	hub          *ws.Hub
	router       *gin.Engine
	guildService usecase.GuildServiceContract
	token        *token.JWTService
}

func (h *WebSocketHandler) Register() {
	basePath := "/api/v1/ws"

	v1 := h.router.Group(basePath)
	{
		v1.POST("/guilds/:guild_id/channels", middleware.AuthMiddleware(h.token), h.createChannelEndpoint)
	}
}
