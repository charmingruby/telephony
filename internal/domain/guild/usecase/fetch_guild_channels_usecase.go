package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) FetchGuildChannels(dto dto.FetchGuildChannelsDTO) ([]entity.Channel, error) {
	_, err  := s.guildRepo.FindByID(dto.GuildID)
	if err != nil {
		return nil, validation.NewNotFoundErr("guild")
	}

	channels, err := s.channelRepo.ListChannelsByGuildID(dto.GuildID, dto.Pagination.Page)
	if err != nil {
		return nil, err
	}

	return channels, nil
}
