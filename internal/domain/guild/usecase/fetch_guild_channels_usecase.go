package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) FetchGuildChannels(dto dto.FetchGuildChannelsDTO) ([]entity.Channel, error) {
	if err := s.userProfileValidation(
		dto.ProfileID,
		dto.UserID,
	); err != nil {
		return nil, err
	}

	_, err := s.guildRepo.FindByID(dto.GuildID)
	if err != nil {
		return nil, validation.NewNotFoundErr("guild")
	}

	if _, err := s.memberRepo.IsAGuildMember(dto.ProfileID, dto.UserID, dto.GuildID); err != nil {
		return nil, validation.NewUnauthorizedErr()
	}

	channels, err := s.channelRepo.ListChannelsByGuildID(dto.GuildID, dto.Pagination.Page)
	if err != nil {
		return nil, err
	}

	return channels, nil
}
