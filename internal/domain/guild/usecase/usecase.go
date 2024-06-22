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
	FetchGuilds(pagination core.PaginationParams) ([]entity.Guild, error)
	DeleteGuild(guildID int) error
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
