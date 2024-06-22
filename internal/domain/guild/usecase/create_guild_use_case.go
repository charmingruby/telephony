package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) CreateGuild(dto dto.CreateGuildDTO) error {
	if err := s.userCliet.UserExists(dto.OwnerID); err != nil {
		return validation.NewNotFoundErr("user")
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

	if err := s.guildRepo.Store(guild); err != nil {
		return validation.NewInternalErr()
	}

	return nil
}
