package entity

import (
	"time"

	"github.com/charmingruby/telephony/internal/core"
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

	// validation

	return &user, nil
}

type User struct {
	ID           int        `json:"id" db:"id"`
	FirstName    string     `json:"first_name" db:"first_name"`
	LastName     string     `json:"last_name" db:"last_name"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}
