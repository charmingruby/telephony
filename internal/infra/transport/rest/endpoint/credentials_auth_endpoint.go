package endpoint

import (
	"github.com/charmingruby/telephony/internal/domain/user/dto"
	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

type CredentialsAuthRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CredentialsAuthData struct {
	AccessToken string `json:"access_token"`
}

type CredentialsAuthResponse struct {
	Message string              `json:"message"`
	Data    CredentialsAuthData `json:"data"`
	Code    int                 `json:"status_code"`
}

// Auth with credentials godoc
//
//	@Summary		Authenticates an user
//	@Description	Authenticates an user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CredentialsAuthRequest	true	"Credentials Payload"
//	@Success		201		{object}	CredentialsAuthResponse
//	@Failure		422		{object}	Response
//	@Failure		409		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/auth/login [post]
func (h *Handler) credentialsAuthEndpoint(c *gin.Context) {
	var req CredentialsAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		connhelper.NewPayloadError(c, err)
		return
	}

	dto := dto.CredentialsAuthDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := h.userService.CredentialsAuth(dto)
	if err != nil {
		invalidCredentialsErr, ok := err.(*validation.ErrInvalidCredentials)
		if ok {
			connhelper.NewInvalidCredentialsError(c, invalidCredentialsErr)
			return
		}

		connhelper.NewInternalServerError(c, err)
		return
	}

	token, err := h.token.GenerateToken(res.UserID)
	if err != nil {
		connhelper.NewInternalServerError(c, err)
		return
	}

	data := CredentialsAuthData{
		AccessToken: token,
	}

	connhelper.NewOkResponse(c, "user authenticated successfully", data)
}
