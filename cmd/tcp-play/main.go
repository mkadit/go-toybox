package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mkadit/go-toybox/internal/adapters/primary/tcp_server"
	"github.com/mkadit/go-toybox/internal/adapters/secondary/tcp_client"
	"github.com/mkadit/go-toybox/internal/applications/service"
	"github.com/mkadit/go-toybox/internal/models"
)

func ProcessUserData(data []byte) ([]byte, error) {
	return append(data, []byte(" processed")...), nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// Initialize the connection repository
	repo := service.NewConnectionRepository()
	defer repo.CloseAllConnections()

	// Load configuration
	config := models.TCPConfiguration{
		Role:       "server", // Change to "server", "client", or "middleware"
		HostIp:     "127.0.0.1",
		HostPort:   8082,
		ClientIp:   "127.0.0.1",
		ClientPort: 8082,
		Timeout:    40,
		PoolSize:   10,
	}

	switch config.Role {
	case "server":
		handler := service.NewConnectionHandler(ctx, repo, ProcessUserData)
		server := tcp_server.NewTCPServer(handler)
		go func() { server.Start(config) }()

	case "client":
		client := tcp_client.NewTCPClient(fmt.Sprintf("%s:%d", config.ClientIp, config.ClientPort))
		clientHandler := service.NewClientHandler(ctx, client, ProcessUserData)
		clientHandler.ConnectAndSend([]byte("Hello from client"))

	case "middleware":
		client := tcp_client.NewTCPClient(fmt.Sprintf("%s:%d", config.ClientIp, config.ClientPort))
		if err := client.Connect(); err != nil {
			fmt.Println("Failed to connect to downstream server:", err)
			return
		}
		defer client.Close()

		middlewareHandler := service.NewMiddlewareHandler(ctx, repo, client, ProcessUserData)
		server := tcp_server.NewTCPServer(middlewareHandler)
		go func() { server.Start(config) }()
	}

	<-sigs
	fmt.Println("Shutting down application...")
}
