package ws

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	connections map[uuid.UUID]*websocket.Conn
	mu          sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uuid.UUID]*websocket.Conn),
	}
}

func (m *ConnectionManager) Add(
	userID uuid.UUID,
	conn *websocket.Conn,
) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.connections[userID] = conn
}

func (m *ConnectionManager) Remove(userID uuid.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.connections, userID)
}

func (m *ConnectionManager) Get(
	userID uuid.UUID,
) (*websocket.Conn, bool) {

	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, exists := m.connections[userID]

	return conn, exists
}