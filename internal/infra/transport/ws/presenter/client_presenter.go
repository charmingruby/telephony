package presenter

import "github.com/charmingruby/telephony/internal/infra/transport/ws"

type ClientHTTP struct {
	ProfileID   int    `json:"profile_id"`
	RoomID      int    `json:"room_id"`
	GuildID     int    `json:"guild_id"`
	DisplayName string `json:"display_name"`
}

func ToHTTP(c ws.Client) ClientHTTP {
	return ClientHTTP{
		ProfileID:   c.ProfileID,
		RoomID:      c.RoomID,
		GuildID:     c.GuildID,
		DisplayName: c.DisplayName,
	}
}
