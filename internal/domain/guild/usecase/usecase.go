package usecase

import (
	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/domain/guild/repository"
	"github.com/charmingruby/telephony/internal/shared/client"
)

type GuildServiceContract interface {
	CreateGuild(dto dto.CreateGuildDTO) error
	FetchAvailableGuilds(pagination core.PaginationParams) ([]entity.Guild, error)
	// DeleteGuild(dto dto.DeleteGuildDTO) error
}

func NewGuildService(
	guildRepo repository.GuildRepository,
	userCliet client.UserClient,
) *GuildService {
	return &GuildService{
		guildRepo: guildRepo,
		userCliet: userCliet,
	}
}

type GuildService struct {
	guildRepo repository.GuildRepository
	userCliet client.UserClient
}