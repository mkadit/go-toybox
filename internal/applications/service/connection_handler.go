package service

import (
	"context"
	"fmt"
	"net"
)

// ConnectionHandler handles incoming connections for the server role.
type ConnectionHandler struct {
	ctx         context.Context
	repo        *ConnectionRepository
	processData func([]byte) ([]byte, error)
}

func NewConnectionHandler(ctx context.Context, repo *ConnectionRepository, processData func([]byte) ([]byte, error)) *ConnectionHandler {
	return &ConnectionHandler{
		ctx:         ctx,
		repo:        repo,
		processData: processData,
	}
}

// HandleConnection implements ports.TCPHandler.
func (h *ConnectionHandler) HandleConnection(conn net.Conn) {
	defer conn.Close()
	id := conn.RemoteAddr().String()
	h.repo.AddConnection(id, conn)
	defer h.repo.RemoveConnection(id)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		data, err := h.processData(buffer[:n])
		fmt.Println(data)
		if err != nil {
			fmt.Println("Error processing data:", err)
			return
		}

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
	}
}
