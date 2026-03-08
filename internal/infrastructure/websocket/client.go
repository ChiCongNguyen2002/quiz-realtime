package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn      *websocket.Conn
	Send      chan []byte
	Hub       *Hub
	SessionID string
}

func NewClient(conn *websocket.Conn, hub *Hub, sessionID string) *Client {
	return &Client{
		Conn:      conn,
		Send:      make(chan []byte, 256),
		Hub:       hub,
		SessionID: sessionID,
	}
}

func (c *Client) WritePump() {
	defer func() {
		_ = c.Conn.Close()
	}()

	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("write message error:", err)
			return
		}
	}
}
