package websocket

import (
	"sync"

	"github.com/chat-socio/backend/pkg/uuid"
	"github.com/hertz-contrib/websocket"
)

type WebSocket struct {
	mapConnections map[string]*WSConnection
	lock           *sync.RWMutex
}

func NewWebSocket() *WebSocket {
	return &WebSocket{
		mapConnections: make(map[string]*WSConnection),
		lock:           &sync.RWMutex{},
	}
}
func (ws *WebSocket) AddConnection(conn *websocket.Conn) (*WSConnection, error) {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	id, err := uuid.NewID()
	if err != nil {
		return nil, err
	}

	wsConn := &WSConnection{
		conn: conn,
		id:   id,
	}

	ws.mapConnections[id] = wsConn
	return wsConn, nil
}

func (ws *WebSocket) AddWrapConnection(wsConn *WSConnection) {
	ws.lock.Lock()
	defer ws.lock.Unlock()

	ws.mapConnections[wsConn.id] = wsConn
}

func (ws *WebSocket) GetConnection(id string) (*WSConnection, bool) {
	ws.lock.RLock()
	defer ws.lock.RUnlock()
	conn, ok := ws.mapConnections[id]
	return conn, ok
}
func (ws *WebSocket) RemoveConnection(id string) {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	if conn, ok := ws.mapConnections[id]; ok {
		conn.Close()
		delete(ws.mapConnections, id)
	}
}
func (ws *WebSocket) GetAllConnections() []*WSConnection {
	ws.lock.RLock()
	defer ws.lock.RUnlock()
	connections := make([]*WSConnection, 0, len(ws.mapConnections))
	for _, conn := range ws.mapConnections {
		connections = append(connections, conn)
	}
	return connections
}
