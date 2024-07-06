package ws

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

type Room struct {
	Clients map[string]*Client `json:"clients"`
	ID      int                `json:"id"`
	Name    string             `json:"name"`
	GuildID int                `json:"guild_id"`
}

type Hub struct {
	Rooms map[string]*Room
}
