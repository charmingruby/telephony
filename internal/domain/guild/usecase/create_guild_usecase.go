package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) CreateGuild(dto dto.CreateGuildDTO) error {
	if err := s.userProfileValidation(
		dto.ProfileID,
		dto.UserID,
	); err != nil {
		return err
	}

	if _, err := s.guildRepo.FindByName(dto.Name); err == nil {
		return validation.NewConflictErr("guild", "name")
	}

	guild, err := entity.NewGuild(
		dto.Name,
		dto.Description,
		dto.ProfileID,
	)
	if err != nil {
		return err
	}

	_, err = s.guildRepo.Store(guild)
	if err != nil {
		return validation.NewInternalErr()
	}

	if err := s.userClient.GuildJoin(dto.ProfileID); err != nil {
		return validation.NewInternalErr()
	}

	return nil
}
