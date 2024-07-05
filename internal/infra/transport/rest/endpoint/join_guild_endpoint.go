package endpoint

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/gin-gonic/gin"
)

type JoinGuildRequest struct {
	ProfileID int `json:"profile_id" binding:"required"`
}

// Create guild godoc
//
//	@Summary		Join a guild
//	@Description	Join a guild
//	@Tags			Members
//	@Accept			json
//	@Produce		json
//	@Param			request	body		JoinGuildRequest	true	"Join Guild Payload"
//	@Success		200		{object}	Response
//	@Failure		400		{object}	Response
//	@Failure		401		{object}	Response
//	@Failure		404		{object}	Response
//	@Failure		409		{object}	Response
//	@Failure		422		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/guilds/{guild_id}/join [post]
func (h *Handler) joinGuildEndpoint(c *gin.Context) {
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

	var req JoinGuildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		connhelper.NewPayloadError(c, err)
		return
	}

	dto := dto.JoinGuildDTO{
		GuildID:   guildID,
		ProfileID: req.ProfileID,
		UserID:    userID,
	}

	if err := h.guildService.JoinGuild(dto); err != nil {
		notFoundErr, ok := err.(*validation.ErrNotFound)
		if ok {
			connhelper.NewResourceNotFoundError(c, notFoundErr)
			return
		}

		badRequestErr, ok := err.(*validation.ErrBadRequest)
		if ok {
			connhelper.NewBadRequestError(c, badRequestErr)
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

	connhelper.NewOkResponse(c, "user joined to the guild successfully", nil)
}
