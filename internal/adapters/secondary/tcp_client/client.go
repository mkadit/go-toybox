package tcp_client

import (
	"fmt"
	"net"
)

// TCPClient represents a TCP client.
type TCPClient struct {
	address string
	conn    net.Conn
}

// NewTCPClient creates a new TCP client.
func NewTCPClient(address string) *TCPClient {
	return &TCPClient{
		address: address,
	}
}

// Connect establishes a connection to the server.
func (c *TCPClient) Connect() error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// Send sends data to the server and waits for a response.
func (c *TCPClient) Send(data []byte) ([]byte, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("connection not established")
	}

	_, err := c.conn.Write(data)
	if err != nil {
		return nil, err
	}

	// Read response
	buffer := make([]byte, 1024)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

// Close closes the connection.
func (c *TCPClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
