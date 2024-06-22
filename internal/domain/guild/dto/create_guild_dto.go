package dto

type CreateGuildDTO struct {
	Name        string
	Description string
	Tags        []string
	OwnerID     int
}
