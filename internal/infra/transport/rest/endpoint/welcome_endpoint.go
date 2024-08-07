package endpoint

import (
	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/gin-gonic/gin"
)

// Welcome godoc
//
//	@Summary		Health Check
//	@Description	Health Check
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	Response
//	@Router			/welcome [get]
func welcomeEndpoint(c *gin.Context) {
	connhelper.NewOkResponse(c, "OK!", nil)
}
