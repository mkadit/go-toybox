package service

import (
	"context"
	"fmt"
	"net"

	"github.com/mkadit/go-toybox/internal/ports"
)

// MiddlewareHandler handles incoming connections and forwards data to a downstream client.
type MiddlewareHandler struct {
	ctx         context.Context
	repo        *ConnectionRepository
	client      ports.TCPClient
	processData func([]byte) ([]byte, error)
}

// NewMiddlewareHandler creates a new middleware handler.
func NewMiddlewareHandler(ctx context.Context, repo *ConnectionRepository, client ports.TCPClient, processData func([]byte) ([]byte, error)) *MiddlewareHandler {
	return &MiddlewareHandler{
		ctx:         ctx,
		repo:        repo,
		client:      client,
		processData: processData,
	}
}

// HandleConnection handles incoming connections from clients.
func (h *MiddlewareHandler) HandleConnection(conn net.Conn) {
	defer conn.Close()
	id := conn.RemoteAddr().String()
	h.repo.AddConnection(id, conn)
	defer h.repo.RemoveConnection(id)

	buffer := make([]byte, 1024)
	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		// Process the data
		processedData, err := h.processData(buffer[:n])
		fmt.Println(processedData)
		if err != nil {
			fmt.Println("Error processing data:", err)
			return
		}

		// Send the processed data to the downstream server
		serverResponse, err := h.client.Send(processedData)
		if err != nil {
			fmt.Println("Error forwarding data to downstream server:", err)
			return
		}

		// Send the server's response back to the client
		_, err = conn.Write(serverResponse)
		if err != nil {
			fmt.Println("Error sending response back to client:", err)
			return
		}
	}
}
