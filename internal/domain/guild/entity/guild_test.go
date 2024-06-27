package entity

import (
	"testing"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/stretchr/testify/assert"
)

func Test_NewGuild(t *testing.T) {
	name := "dummy name"
	description := "dummy description"
	tags := []string{"tag1"}
	profileID := 2

	t.Run("it should be able to create a guild", func(t *testing.T) {
		g, err := NewGuild(name, description, tags, profileID)

		assert.NoError(t, err)
		assert.Equal(t, core.NewDefaultDomainID(), g.ID)
		assert.Equal(t, name, g.Name)
		assert.Equal(t, description, g.Description)
		assert.Equal(t, tags, g.Tags)
		assert.Equal(t, 0, g.ChannelsQuantity)
		assert.Equal(t, profileID, g.OwnerID)
		assert.Nil(t, g.DeletedAt)
	})

	t.Run("it should receive an error when try to create a guild with invalid params", func(t *testing.T) {
		g, err := NewGuild("", description, tags, profileID)

		assert.Nil(t, g)
		assert.Error(t, err)
		assert.Equal(t, validation.NewValidationErr(validation.ErrRequired("name")).Error(), err.Error())
	})
}
