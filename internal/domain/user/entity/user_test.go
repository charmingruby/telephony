package entity

import (
	"testing"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/validation"
	"github.com/stretchr/testify/assert"
)

func Test_NewUser(t *testing.T) {
	firstName := "john"
	lastName := "doe"
	email := "john@doe.com"
	password := "password123"

	t.Run("it should be able to create an user", func(t *testing.T) {
		u, err := NewUser(
			firstName,
			lastName,
			email,
			password,
		)

		assert.NoError(t, err)
		assert.Equal(t, core.NewDefaultDomainID(), u.ID)
		assert.Equal(t, firstName, u.FirstName)
		assert.Equal(t, lastName, u.LastName)
		assert.Equal(t, email, u.Email)
		assert.Equal(t, password, u.PasswordHash)
		assert.Nil(t, u.DeletedAt)
	})

	t.Run("it should receive an error when try to create an user with invalid params", func(t *testing.T) {
		u, err := NewUser(
			"",
			lastName,
			email,
			password,
		)

		assert.Nil(t, u)
		assert.Error(t, err)
		assert.Equal(t, validation.NewValidationErr(validation.ErrRequired("firstname")).Error(), err.Error())
	})
}
