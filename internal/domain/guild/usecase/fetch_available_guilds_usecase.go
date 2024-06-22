package usecase

import (
	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) FetchAvailableGuilds(pagination core.PaginationParams) ([]entity.Guild, error) {
	guilds, err := s.guildRepo.ListAvailables(pagination.Page)
	if err != nil {
		return []entity.Guild{}, validation.NewInternalErr()
	}

	return guilds, nil
}
