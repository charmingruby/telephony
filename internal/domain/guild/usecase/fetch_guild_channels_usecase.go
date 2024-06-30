package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) FetchGuildChannels(dto dto.FetchGuildChannelsDTO) ([]entity.Channel, error) {
	userExists := s.userClient.UserExists(dto.UserID)
	if !userExists {
		return nil, validation.NewNotFoundErr("user")
	}

	profileExists := s.userClient.ProfileExists(dto.ProfileID)
	if !profileExists {
		return nil, validation.NewNotFoundErr("user_profile")
	}

	isTheProfileOwner := s.userClient.IsTheProfileOwner(dto.UserID, dto.ProfileID)
	if !isTheProfileOwner {
		return nil, validation.NewUnauthorizedErr()
	}

	_, err := s.guildRepo.FindByID(dto.GuildID)
	if err != nil {
		return nil, validation.NewNotFoundErr("guild")
	}

	channels, err := s.channelRepo.ListChannelsByGuildID(dto.GuildID, dto.Pagination.Page)
	if err != nil {
		return nil, err
	}

	return channels, nil
}
