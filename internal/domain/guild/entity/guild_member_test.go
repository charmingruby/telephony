package entity

import (
	"testing"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/stretchr/testify/assert"
)

func Test_NewGuildMember(t *testing.T) {
	profileID := 1
	userID := 1
	guildID := 1

	t.Run("it should be able to create a guild member", func(t *testing.T) {
		m, err := NewGuildMember(profileID, userID, guildID)

		assert.NoError(t, err)
		assert.Equal(t, core.NewDefaultDomainID(), m.ID)
		assert.Equal(t, profileID, m.ProfileID)
		assert.Equal(t, userID, m.UserID)
		assert.Equal(t, guildID, m.GuildID)
	})
}
