package dto

type CreateChannelDTO struct {
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}
