package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) CreateChannel(dto dto.CreateChannelDTO) error {
	if err := s.userProfileValidation(
		dto.ProfileID,
		dto.UserID,
	); err != nil {
		return err
	}

	guild, err := s.guildRepo.FindByID(dto.GuildID)
	if err != nil {
		return validation.NewNotFoundErr("guild")
	}

	if _, err := s.memberRepo.IsAGuildMember(dto.ProfileID, dto.UserID, dto.GuildID); err != nil {
		return validation.NewUnauthorizedErr()
	}

	if guild.OwnerID != dto.ProfileID {
		return validation.NewUnauthorizedErr()
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

	if err := s.userClient.GuildJoin(dto.ProfileID); err != nil {
		return err
	}

	return nil
}
