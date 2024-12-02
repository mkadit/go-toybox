package ports

import "github.com/mkadit/go-toybox/internal/models"

type DbPort interface {
	CreateUser(user models.User) (err error)
	GetUserByEmail(email string) (result models.User, err error)
}

// TCPClient defines the interface for TCP client operations.
type TCPClient interface {
	Connect() error
	Send(data []byte) ([]byte, error) // Send returns a response and an error
	Close() error                     // Close returns an error
}
