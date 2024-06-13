package core

import "github.com/google/uuid"

func NewUUID() string {
	id, _ := uuid.NewUUID()
	return id.String()
}
