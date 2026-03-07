package http

import (
	"net/http"

	"github.com/gorilla/websocket"
	ws "quiz-realtime/internal/infrastructure/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(hub *ws.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := ws.NewClient(conn, hub)
	hub.Register <- client

	go client.WritePump()
}

