package entity

import (
	"testing"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/stretchr/testify/assert"
)

func Test_NewChannel(t *testing.T) {
	name := "dummy name"
	guildID := 1
	profileID := 2

	t.Run("it should be able to create a channel", func(t *testing.T) {
		c, err := NewChannel(name, guildID, profileID)

		assert.NoError(t, err)
		assert.Equal(t, core.NewDefaultDomainID(), c.ID)
		assert.Equal(t, name, c.Name)
		assert.Equal(t, 0, c.MessagesQuantity)
		assert.Equal(t, guildID, c.GuildID)
		assert.Equal(t, profileID, c.OwnerID)
		assert.Nil(t, c.DeletedAt)
	})

	t.Run("it should receive an error when try to create a channel with invalid params", func(t *testing.T) {
		c, err := NewChannel("", guildID, profileID)

		assert.Nil(t, c)
		assert.Error(t, err)
		assert.Equal(t, validation.NewValidationErr(validation.ErrRequired("name")).Error(), err.Error())
	})
}
