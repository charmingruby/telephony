package repository

import "github.com/charmingruby/telephony/internal/domain/guild/entity"

type GuildMemberRepository interface {
	Store(m *entity.GuildMember) (int, error)
	IsAGuildMember(profileID, userID, guildID int) (*entity.GuildMember, error)
}
