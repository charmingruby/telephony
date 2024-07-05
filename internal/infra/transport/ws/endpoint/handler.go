package endpoint

import (
	"github.com/charmingruby/telephony/internal/infra/transport/common/middleware"
	"github.com/charmingruby/telephony/internal/infra/transport/ws"
)

func NewWSHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

type WebSocketHandler struct {
	hub *ws.Hub
}

func (ws *WebSocketHandler) Register() {
	basePath := "/api/v1/ws"

	v1 := ws.hub.Router.Group(basePath)
	{
		v1.POST("/guilds/:guild_id/channels", middleware.AuthMiddleware(ws.hub.Token), ws.createChannelEndpoint)
	}
}
