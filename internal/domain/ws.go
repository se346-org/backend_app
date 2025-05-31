package domain

import "github.com/chat-socio/backend/infrastructure/websocket"

var WebSocket = &websocket.WebSocket{}

func InitWebSocket() {
	WebSocket = websocket.NewWebSocket()
}
