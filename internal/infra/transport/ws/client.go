package ws

import "github.com/gorilla/websocket"

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       int    `json:"id"`
	RoomID   string `json:"room_id"`
	SenderID int    `json:"sender_id"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"room_id"`
	SenderID int    `json:"sender_id"`
}
