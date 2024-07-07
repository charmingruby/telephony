package ws

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewClient(
	c *gin.Context,
	hub *Hub,
	displayName string, profileID, channelID, guildID int) (*Client, error) {
	conn, err := hub.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}

	cl := Client{
		Conn:        conn,
		MessageCh:   make(chan *Message, 10),
		ProfileID:   profileID,
		RoomID:      channelID,
		GuildID:     guildID,
		DisplayName: displayName,
	}

	joinMsg := Message{
		Content: fmt.Sprintf("A wild %s has appeared", displayName),
		RoomID:  channelID,
		GuildID: guildID,
	}

	hub.RegisterCh <- &cl
	hub.BroadcastCh <- &joinMsg

	return &cl, nil
}

type Client struct {
	Conn        *websocket.Conn
	MessageCh   chan *Message
	ProfileID   int    `json:"profile_id"`
	RoomID      int    `json:"room_id"`
	GuildID     int    `json:"guild_id"`
	DisplayName string `json:"display_name"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   int    `json:"room_id"`
	GuildID  int    `json:"guild_id"`
	SenderID int    `json:"sender_id"`
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		msg, ok := <-c.MessageCh
		if !ok {
			return
		}

		c.Conn.WriteJSON(msg)
	}
}

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.UnregisterCh <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error(fmt.Sprintf("[WEBSOCKET] Error: %v", err))
			}
			break
		}

		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			GuildID:  c.GuildID,
			SenderID: c.ProfileID,
		}

		hub.BroadcastCh <- msg
	}
}
