package websocket

import (
	"github.com/chat-socio/backend/pkg/uuid"
	"github.com/hertz-contrib/websocket"
)

type WSConnection struct {
	conn *websocket.Conn
	id   string
}

func (wsConn *WSConnection) SendMessage(message []byte) error {
	err := wsConn.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (wsConn *WSConnection) ReceiveMessage() ([]byte, error) {
	_, message, err := wsConn.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (wsConn *WSConnection) Close() error {
	err := wsConn.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (wsConn *WSConnection) GetID() string {
	return wsConn.id
}

func (wsConn *WSConnection) GetConn() *websocket.Conn {
	return wsConn.conn
}

func (wsConn *WSConnection) SetConn(conn *websocket.Conn) {
	wsConn.conn = conn
}

func (wsConn *WSConnection) SetID(id string) {
	wsConn.id = id
}

func NewWSConnection(conn *websocket.Conn) (*WSConnection, error) {
	// Generate a unique ID for the connection
	id, err := uuid.NewID()
	if err != nil {
		return nil, err
	}
	return &WSConnection{
		conn: conn,
		id:   id,
	}, nil
}
