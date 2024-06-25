package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) CreateGuild(dto dto.CreateGuildDTO) error {
	profileExists := s.userCliet.ProfileExists(dto.OwnerID)
	if !profileExists {
		return validation.NewNotFoundErr("user_profile")
	}

	guild, err := entity.NewGuild(
		dto.Name,
		dto.Description,
		dto.Tags,
		dto.OwnerID,
	)
	if err != nil {
		return err
	}

	_, err = s.guildRepo.Store(guild)
	if err != nil {
		return validation.NewInternalErr()
	}

	if err := s.userCliet.GuildJoin(dto.OwnerID); err != nil {
		return validation.NewInternalErr()
	}

	return nil
}
