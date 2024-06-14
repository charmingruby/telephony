package dto

type CreateProfileDTO struct {
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	UserID      int    `json:"user_id"`
}
