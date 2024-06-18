package entity

import (
	"time"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/validation"
)

func NewUser(firstName, lastName, email, password string) (*User, error) {
	user := User{
		ID:           core.NewDefaultDomainID(),
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: password, // will be hashed on creation
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		DeletedAt:    nil,
	}

	if err := validation.ValidateStruct(user); err != nil {
		return nil, err
	}

	return &user, nil
}

type User struct {
	ID           int        `json:"id" validate:"required" db:"id"`
	FirstName    string     `json:"first_name" validate:"required,min=1,max=36" db:"first_name"`
	LastName     string     `json:"last_name" validate:"required,min=1" db:"last_name"`
	Email        string     `json:"email" validate:"required,email" db:"email"`
	PasswordHash string     `json:"password_hash" validate:"required,min=8,max=24"  db:"password_hash"`
	CreatedAt    time.Time  `json:"created_at" validate:"required" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" validate:"required" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}
