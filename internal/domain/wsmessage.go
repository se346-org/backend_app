package domain

import "encoding/json"

type WebSocketMessageType string

const (
	WsAuthorization     = "AUTHORIZATION"
	WsMessage           = "MESSAGE"
	WsPing              = "PING"
	WsPong              = "PONG"
	WsUpdateLastMessage = "UPDATE_LAST_MESSAGE"
	WsSeenMessage       = "SEEN_MESSAGE"
)

// WebSocketMessage represents a message sent over a WebSocket connection.
// It contains a type and a payload, where the payload can be of any type.
// The type is a string that indicates the kind of message being sent.
// The payload is a generic type that can be any data structure.
// The WebSocketMessage struct implements the json.Marshaler and json.Unmarshaler interfaces,
// allowing it to be easily converted to and from JSON format.
type WebSocketMessage struct {
	Type              WebSocketMessageType `json:"type,omitempty"`
	Payload           map[string]any       `json:"payload,omitempty"`
	IgnoreUserOnlines []string             `json:"ignore_user_onlines,omitempty"`
}

func NewWebSocketMessage(messageType WebSocketMessageType, payload map[string]any) *WebSocketMessage {
	return &WebSocketMessage{
		Type:    messageType,
		Payload: payload,
	}
}

func (wsm *WebSocketMessage) GetType() WebSocketMessageType {
	return wsm.Type
}

func (wsm *WebSocketMessage) GetPayload() map[string]any {
	return wsm.Payload
}

func (wsm *WebSocketMessage) SetType(messageType WebSocketMessageType) {
	wsm.Type = messageType
}

func (wsm *WebSocketMessage) SetPayload(payload map[string]any) {
	wsm.Payload = payload
}

func (wsm *WebSocketMessage) String() string {
	data, err := json.Marshal(wsm)
	if err != nil {
		return ""
	}
	return string(data)
}
