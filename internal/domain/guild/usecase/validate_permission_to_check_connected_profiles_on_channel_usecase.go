package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) ValidatePermissionToCheckConnectedProfilesOnChannel(
	dto dto.ValidatePermissionToCheckConnectedProfilesOnChannelDTO,
) error {
	if err := s.userProfileValidation(
		dto.ProfileID,
		dto.UserID,
	); err != nil {
		return err
	}

	_, err := s.guildRepo.FindByID(dto.GuildID)
	if err != nil {
		return validation.NewNotFoundErr("guild")
	}

	if _, err := s.memberRepo.IsAGuildMember(dto.ProfileID, dto.UserID, dto.GuildID); err != nil {
		return validation.NewUnauthorizedErr()
	}

	return nil
}
