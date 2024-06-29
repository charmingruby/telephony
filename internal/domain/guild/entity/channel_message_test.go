package entity

import (
	"testing"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/stretchr/testify/assert"
)

func Test_NewChannelMessage(t *testing.T) {
	content := "hello world"
	senderID := 2

	t.Run("it should be able to create a channel message", func(t *testing.T) {
		c, err := NewChannelMessage(content, senderID)

		assert.NoError(t, err)
		assert.Equal(t, core.NewDefaultDomainID(), c.ID)
		assert.Equal(t, content, c.Content)
		assert.Equal(t, senderID, c.SenderID)
		assert.Nil(t, c.DeletedAt)
	})

	t.Run("it should receive an error when try to create a channel message with invalid params", func(t *testing.T) {
		c, err := NewChannelMessage("", senderID)

		assert.Nil(t, c)
		assert.Error(t, err)
		assert.Equal(t, validation.NewValidationErr(validation.ErrRequired("content")).Error(), err.Error())
	})
}
