package api

import (
	"github.com/mkadit/go-toybox/internal/models"
	"github.com/mkadit/go-toybox/internal/ports"
)

type (
	Users interface {
		HashPassword(password string) (string, error)
		CheckPasswordHash(password, hash string) bool
	}
	Urls interface{}
)

// Application implements the APIPort interface
type Application struct {
	db          ports.DbPort
	users       Users
	urls        Urls
	emailConfig models.EmailConfiguration
}

// NewApplication Application layer to connect domains and adapters
func NewApplication(db ports.DbPort, users Users, urls Urls, emailConfig models.EmailConfiguration) *Application {
	return &Application{db: db, users: users, urls: urls, emailConfig: emailConfig}
}
