package usecase

import (
	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/domain/user/adapter"
	"github.com/charmingruby/telephony/internal/domain/user/repository"
)

type GuildServiceContract interface {
	CreateGuild(dto dto.CreateGuildDTO) error
	FetchGuilds(pagination core.PaginationParams) ([]entity.Guild, error)
	DeleteGuild(guildID int) error
}

func NewGuildService(
	userRepo repository.UserRepository,
	profileRepo repository.UserProfileRepository,
	crypto adapter.CryptographyContract,
) *GuildService {
	return &GuildService{}
}

type GuildService struct{}
