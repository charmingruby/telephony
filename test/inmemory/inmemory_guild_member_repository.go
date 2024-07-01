package inmemory

import (
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func NewInMemoryGuildMemberRepository() *InMemoryGuildMemberRepository {
	return &InMemoryGuildMemberRepository{
		Items: []entity.GuildMember{},
	}
}

type InMemoryGuildMemberRepository struct {
	Items []entity.GuildMember
}

func (r *InMemoryGuildMemberRepository) Store(m *entity.GuildMember) (int, error) {
	r.Items = append(r.Items, *m)
	return m.ID, nil
}

func (r *InMemoryGuildMemberRepository) IsAGuildMember(profileID, userID, guildID int) (*entity.GuildMember, error) {
	for _, e := range r.Items {
		if e.ProfileID == profileID && e.GuildID == guildID && e.UserID == userID && e.IsActive {
			return &e, nil
		}
	}

	return nil, validation.NewNotFoundErr("guild member")
}
