package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) JoinChannel(pl dto.JoinChannelDTO) (*dto.JoinChannelResponseDTO, error) {
	if err := s.userProfileValidation(
		pl.ProfileID,
		pl.UserID,
	); err != nil {
		return nil, err
	}

	_, err := s.guildRepo.FindByID(pl.GuildID)
	if err != nil {
		return nil, validation.NewNotFoundErr("guild")
	}

	if _, err := s.memberRepo.IsAGuildMember(pl.ProfileID, pl.UserID, pl.GuildID); err != nil {
		return nil, validation.NewUnauthorizedErr()
	}

	if _, err := s.channelRepo.FindByID(pl.ChannelID, pl.GuildID); err != nil {
		return nil, validation.NewNotFoundErr("channel")
	}

	displayName, err := s.userClient.GetDisplayName(pl.ProfileID)
	if err != nil {
		return nil, validation.NewNotFoundErr("profile")
	}

	res := dto.JoinChannelResponseDTO{DisplayName: displayName}

	return &res, nil
}
