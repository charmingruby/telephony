package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) JoinGuild(dto dto.JoinGuildDTO) error {
	if err := s.userProfileValidation(dto.ProfileID, dto.UserID); err != nil {
		return err
	}

	if _, err := s.memberRepo.IsAGuildMember(dto.ProfileID, dto.UserID, dto.GuildID); err == nil {
		return validation.NewBadRequestErr("is already member")
	}

	member, err := entity.NewGuildMember(dto.ProfileID, dto.UserID, dto.GuildID)
	if err != nil {
		return err
	}

	if _, err := s.memberRepo.Store(member); err != nil {
		return validation.NewInternalErr()
	}

	return nil
}
