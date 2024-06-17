package entity

import (
	"testing"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/stretchr/testify/assert"
)

func Test_NewUserProfile(t *testing.T) {
	displayName := "dummy nick"
	bio := "dummy bio"
	userID := 1

	t.Run("it should be able to create a profile", func(t *testing.T) {
		u, err := NewUserProfile(
			displayName,
			bio,
			userID,
		)

		assert.NoError(t, err)
		assert.Equal(t, core.NewDefaultDomainID(), u.ID)
		assert.Equal(t, displayName, u.DisplayName)
		assert.Equal(t, bio, u.Bio)
		assert.Equal(t, userID, u.UserID)
		assert.Equal(t, 0, u.GuildsQuantity)
		assert.Equal(t, 0, u.MessagesQuantity)
		assert.Nil(t, u.DeletedAt)
	})

	t.Run("it should receive an error when try to create a profile with invalid params", func(t *testing.T) {
		u, err := NewUserProfile(
			"",
			bio,
			userID,
		)

		assert.Nil(t, u)
		assert.Error(t, err)
		assert.Equal(t, validation.NewValidationErr(validation.ErrRequired("displayname")).Error(), err.Error())
	})
}
