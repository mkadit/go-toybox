package api

import "github.com/mkadit/go-toybox/internal/models"

// Register implements ports.APIPort.
func (api *Application) Register(request models.UserRequest) (resp models.GenericResponse, appErr models.ApplicationErr) {
	password, err := api.users.HashPassword(request.Password)
	if err != nil {
		appErr.Wrap(err, models.ErrHashData)
		return
	}
	dbModel := models.User{
		Email:        request.Email,
		Username:     request.Username,
		PasswordHash: password,
	}
	err = api.db.CreateUser(dbModel)
	if err != nil {
		appErr.Wrap(err, models.ErrInsertData)
		return
	}
	return
}
