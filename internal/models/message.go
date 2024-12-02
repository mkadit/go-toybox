package models

type (
	// for key in context
	Request  struct{}
	Response struct{}

	GenericRequest struct {
		Message string `json:"message"`
	}
	GenericResponse struct {
		Message string `json:"message"`
	}
	UserRequest struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=16"`
	}

	LoginResponse struct {
		Token   string `json:"-"`
		Message string `json:"message"`
	}

	EmailData struct {
		Server        string
		Port          int
		EmailSender   string
		EmailPassword string
		Name          string
		EmailReceiver string
		NewPassword   string
	}
)
