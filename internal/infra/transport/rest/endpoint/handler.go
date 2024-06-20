package endpoint

import (
	docs "github.com/charmingruby/telephony/docs"
	exampleUc "github.com/charmingruby/telephony/internal/domain/example/usecase"
	userUc "github.com/charmingruby/telephony/internal/domain/user/usecase"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHandler(router *gin.Engine,
	exampleService exampleUc.ExampleServiceContract,
	userService userUc.UserServiceContract,
) *Handler {
	return &Handler{
		router:         router,
		exampleService: exampleService,
		userService:    userService,
	}
}

type Handler struct {
	router         *gin.Engine
	exampleService exampleUc.ExampleServiceContract
	userService    userUc.UserServiceContract
}

func (h *Handler) Register() {
	basePath := "/api/v1"
	v1 := h.router.Group(basePath)
	docs.SwaggerInfo.BasePath = basePath
	{
		v1.GET("/welcome", welcomeEndpoint)

		// TODO: remove when is finished
		v1.POST("/examples", h.CreateExampleEndpoint)
		v1.GET("/examples/:id", h.getExampleEndpoint)

		v1.POST("/auth/register", h.registerEndpoint)
	}

	h.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
