package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) CreateGuild(dto dto.CreateGuildDTO) error {
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
