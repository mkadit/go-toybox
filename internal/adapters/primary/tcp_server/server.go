package tcp_server

import (
	"fmt"
	"net"

	"github.com/mkadit/go-toybox/internal/models"
	"github.com/mkadit/go-toybox/internal/ports"
)

type TCPServer struct {
	handler ports.TCPHandler
}

func NewTCPServer(handler ports.TCPHandler) *TCPServer {
	return &TCPServer{
		handler: handler,
	}
}

func (s *TCPServer) Start(config models.TCPConfiguration) error {
	hosting := &net.TCPAddr{
		Port: config.HostPort,
	}
	l, err := net.ListenTCP("tcp", hosting)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer l.Close()

	for {
		fmt.Println("Listening TPC on : ", l.Addr().String())
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		go s.handler.HandleConnection(conn)
	}
}
