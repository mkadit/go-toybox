package service

import (
	"net"
	"sync"
)

type ConnectionRepository struct {
	connections sync.Map // Thread-safe map for storing connections
}

// NewConnectionRepository initializes and returns a new ConnectionRepository.
func NewConnectionRepository() *ConnectionRepository {
	return &ConnectionRepository{}
}

// AddConnection stores a connection with a given ID.
func (r *ConnectionRepository) AddConnection(id string, conn net.Conn) {
	r.connections.Store(id, conn)
}

// RemoveConnection deletes a connection by ID and closes it.
func (r *ConnectionRepository) RemoveConnection(id string) {
	if val, ok := r.connections.LoadAndDelete(id); ok {
		conn := val.(net.Conn)
		conn.Close()
	}
}

// GetConnection retrieves a connection by ID.
func (r *ConnectionRepository) GetConnection(id string) (net.Conn, bool) {
	val, ok := r.connections.Load(id)
	if !ok {
		return nil, false
	}
	return val.(net.Conn), true
}

// CloseAllConnections safely closes all connections.
func (r *ConnectionRepository) CloseAllConnections() {
	r.connections.Range(func(key, value any) bool {
		conn := value.(net.Conn)
		conn.Close()
		r.connections.Delete(key)
		return true
	})
}
