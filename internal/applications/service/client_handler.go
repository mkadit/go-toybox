package service

import (
	"context"
	"fmt"

	"github.com/mkadit/go-toybox/internal/adapters/secondary/tcp_client"
)

type ClientHandler struct {
	ctx          context.Context
	client       *tcp_client.TCPClient
	processLogic func(data []byte) ([]byte, error)
}

func NewClientHandler(ctx context.Context, client *tcp_client.TCPClient, processLogic func(data []byte) ([]byte, error)) *ClientHandler {
	return &ClientHandler{ctx: ctx, client: client, processLogic: processLogic}
}

func (h *ClientHandler) ConnectAndSend(data []byte) error {
	if err := h.client.Connect(); err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	defer h.client.Close()

	response, err := h.client.Send(data)
	if err != nil {
		return fmt.Errorf("send error: %w", err)
	}

	processed, err := h.processLogic(response)
	if err != nil {
		return fmt.Errorf("processing error: %w", err)
	}

	fmt.Println("Processed response:", string(processed))
	return nil
}
