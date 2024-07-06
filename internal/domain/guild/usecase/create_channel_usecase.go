package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) CreateChannel(dto dto.CreateChannelDTO) (int, error) {
	if err := s.userProfileValidation(
		dto.ProfileID,
		dto.UserID,
	); err != nil {
		return -1, err
	}

	guild, err := s.guildRepo.FindByID(dto.GuildID)
	if err != nil {
		return -1, validation.NewNotFoundErr("guild")
	}

	if _, err := s.memberRepo.IsAGuildMember(dto.ProfileID, dto.UserID, dto.GuildID); err != nil {
		return -1, validation.NewUnauthorizedErr()
	}

	if guild.OwnerID != dto.ProfileID {
		return -1, validation.NewUnauthorizedErr()
	}

	if _, err := s.channelRepo.FindByName(dto.GuildID, dto.Name); err == nil {
		return -1, validation.NewConflictErr("channel", "name")
	}

	channel, err := entity.NewChannel(
		dto.Name,
		dto.GuildID,
		dto.ProfileID,
	)
	if err != nil {
		return -1, err
	}

	channelID, err := s.channelRepo.Store(channel)
	if err != nil {
		return -1, validation.NewInternalErr()
	}

	if err := s.userClient.GuildJoin(dto.ProfileID); err != nil {
		return -1, err
	}

	return channelID, nil
}
