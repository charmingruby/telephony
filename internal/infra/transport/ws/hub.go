package ws

import (
	"net/http"

	"github.com/charmingruby/telephony/internal/domain/guild/repository"
	"github.com/gorilla/websocket"
)

func NewHub(chRepo repository.ChannelRepository) *Hub {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &Hub{
		Rooms:        make(map[int]*Room),
		Upgrader:     &upgrader,
		RegisterCh:   make(chan *Client),
		UnregisterCh: make(chan *Client),
		BroadcastCh:  make(chan *Message),
		ChannelRepo:  chRepo,
	}
}

func (h *Hub) RegisterRooms() error {
	ch, err := h.ChannelRepo.ListAllChannels()
	if err != nil {
		return err
	}

	for _, c := range ch {
		h.AddRoom(c.ID, c.GuildID, c.Name)
	}

	return nil
}

func (h *Hub) Start() {
	for {
		select {
		case cl := <-h.RegisterCh:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				r := h.Rooms[cl.RoomID]

				if _, ok := r.Clients[cl.ProfileID]; !ok {
					r.Clients[cl.ProfileID] = cl
				}
			}
		case cl := <-h.UnregisterCh:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ProfileID]; !ok {
					delete(h.Rooms[cl.RoomID].Clients, cl.ProfileID)
					close(cl.MessageCh)
				}
			}
		case m := <-h.BroadcastCh:
			if _, ok := h.Rooms[m.RoomID]; ok {
				for _, cl := range h.Rooms[m.RoomID].Clients {
					cl.MessageCh <- m
				}
			}
		}
	}
}

func (h *Hub) AddRoom(channelID, guildID int, name string) {
	h.Rooms[channelID] = &Room{
		ID:      channelID,
		Name:    name,
		GuildID: guildID,
		Clients: make(map[int]*Client),
	}
}

type Room struct {
	Clients map[int]*Client `json:"clients"`
	ID      int             `json:"id"`
	Name    string          `json:"name"`
	GuildID int             `json:"guild_id"`
}

type Hub struct {
	Rooms        map[int]*Room
	Upgrader     *websocket.Upgrader
	RegisterCh   chan *Client
	UnregisterCh chan *Client
	BroadcastCh  chan *Message
	ChannelRepo  repository.ChannelRepository
}
