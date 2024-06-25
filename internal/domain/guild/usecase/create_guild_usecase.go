package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) CreateGuild(dto dto.CreateGuildDTO) error {
	profileExists := s.userCliet.ProfileExists(dto.ProfileID)
	if !profileExists {
		return validation.NewNotFoundErr("user_profile")
	}

	if _, err := s.guildRepo.FindByName(dto.Name); err == nil {
		return validation.NewConflictErr("guild", "name")
	}

	guild, err := entity.NewGuild(
		dto.Name,
		dto.Description,
		dto.Tags,
		dto.ProfileID,
	)
	if err != nil {
		return err
	}

	_, err = s.guildRepo.Store(guild)
	if err != nil {
		return validation.NewInternalErr()
	}

	if err := s.userCliet.GuildJoin(dto.ProfileID); err != nil {
		return validation.NewInternalErr()
	}

	return nil
}
