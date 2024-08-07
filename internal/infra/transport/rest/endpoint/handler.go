package endpoint

import (
	docs "github.com/charmingruby/telephony/docs"
	guildUc "github.com/charmingruby/telephony/internal/domain/guild/usecase"
	userUc "github.com/charmingruby/telephony/internal/domain/user/usecase"

	"github.com/charmingruby/telephony/internal/infra/security/token"
	"github.com/charmingruby/telephony/internal/infra/transport/common/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHandler(
	router *gin.Engine,
	token *token.JWTService,
	userService userUc.UserServiceContract,
	guildService guildUc.GuildServiceContract,
) *Handler {
	return &Handler{
		router:       router,
		token:        token,
		userService:  userService,
		guildService: guildService,
	}
}

type Handler struct {
	router       *gin.Engine
	token        *token.JWTService
	userService  userUc.UserServiceContract
	guildService guildUc.GuildServiceContract
}

func (h *Handler) Register() {
	basePath := "/api/v1"
	docs.SwaggerInfo.BasePath = basePath

	v1 := h.router.Group(basePath)
	{
		v1.GET("/welcome", welcomeEndpoint)
		v1.POST("/auth/register", h.registerEndpoint)
		v1.POST("/auth/login", h.credentialsAuthEndpoint)
		v1.POST("/me/profile", middleware.AuthMiddleware(h.token), h.createProfileEndpoint)
		v1.GET("/me", middleware.AuthMiddleware(h.token), h.meEndpoint)
		v1.POST("/guilds", middleware.AuthMiddleware(h.token), h.createGuildEndpoint)
		v1.GET("/guilds", middleware.AuthMiddleware(h.token), h.fetchAvailableGuildsEndpoint)
		v1.POST("/guilds/:guild_id/join", middleware.AuthMiddleware(h.token), h.joinGuildEndpoint)
	}

	h.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
