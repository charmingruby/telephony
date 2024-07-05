package ws

import (
	"github.com/charmingruby/telephony/internal/domain/guild/usecase"
	"github.com/charmingruby/telephony/internal/infra/security/token"
	"github.com/gin-gonic/gin"
)

func NewHub(
	router *gin.Engine,
	guildService usecase.GuildServiceContract,
	token *token.JWTService,
) *Hub {
	return &Hub{
		Rooms:        make(map[string]*Room),
		Router:       router,
		GuildService: guildService,
		Token:        token,
	}
}

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Router       *gin.Engine
	GuildService usecase.GuildServiceContract
	Token        *token.JWTService
	Rooms        map[string]*Room
}
