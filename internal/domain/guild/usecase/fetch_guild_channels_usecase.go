package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
)

func (s *GuildService) FetchGuildChannels(dto dto.FetchGuildChannelsDTO) ([]entity.Channel, error) {
	return nil, nil
}
