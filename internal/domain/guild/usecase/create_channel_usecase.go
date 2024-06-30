package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) CreateChannel(dto dto.CreateChannelDTO) error {
	userExists := s.userClient.UserExists(dto.UserID)
	if !userExists {
		return validation.NewNotFoundErr("user")
	}

	profileExists := s.userClient.ProfileExists(dto.ProfileID)
	if !profileExists {
		return validation.NewNotFoundErr("user_profile")
	}

	isTheProfileOwner := s.userClient.IsTheProfileOwner(dto.UserID, dto.ProfileID)
	if !isTheProfileOwner {
		return validation.NewUnauthorizedErr()
	}

	if _, err := s.guildRepo.FindByID(dto.GuildID); err != nil {
		return validation.NewNotFoundErr("guild")
	}

	if _, err := s.channelRepo.FindByName(dto.GuildID, dto.Name); err == nil {
		return validation.NewConflictErr("channel", "name")
	}

	channel, err := entity.NewChannel(
		dto.Name,
		dto.GuildID,
		dto.ProfileID,
	)
	if err != nil {
		return err
	}

	if _, err := s.channelRepo.Store(channel); err != nil {
		return validation.NewInternalErr()
	}

	return nil
}
