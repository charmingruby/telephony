package usecase

import (
	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/domain/guild/repository"
	"github.com/charmingruby/telephony/internal/domain/shared/client"
)

type GuildServiceContract interface {
	CreateGuild(dto dto.CreateGuildDTO) error
	FetchAvailableGuilds(pagination core.PaginationParams) ([]entity.Guild, error)
	JoinGuild(dto dto.JoinGuildDTO) error
	CreateChannel(dto dto.CreateChannelDTO) (int, error)
	FetchGuildChannels(dto dto.FetchGuildChannelsDTO) ([]entity.Channel, error)
	JoinChannel(dto dto.JoinChannelDTO) error
	SendMessage(dto dto.SendMessageDTO) error
	DeleteMessage(dto dto.DeleteMessageDTO) error
}

func NewGuildService(
	guildRepo repository.GuildRepository,
	memberRepo repository.GuildMemberRepository,
	channelRepo repository.ChannelRepository,
	userClient client.UserClient,
) *GuildService {
	return &GuildService{
		guildRepo:   guildRepo,
		memberRepo:  memberRepo,
		channelRepo: channelRepo,
		userClient:  userClient,
	}
}

type GuildService struct {
	guildRepo   repository.GuildRepository
	memberRepo  repository.GuildMemberRepository
	channelRepo repository.ChannelRepository
	userClient  client.UserClient
}
