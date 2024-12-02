package ports

import (
	"github.com/mkadit/go-toybox/internal/models"
)

type APIPort interface {
	Register(request models.UserRequest) (models.GenericResponse, models.ApplicationErr)
}
