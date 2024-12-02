package models

import (
	"errors"
)

// General message and error
var (
	ErrAssertData        = errors.New("unable to assert data")
	ErrBadRequest        = errors.New("bad request")
	ErrCreateToken       = errors.New("error token expired")
	ErrCreatingMigration = errors.New("unable create migration object")
	ErrDecryptData       = errors.New("unable to decrypt data")
	ErrDirtyMigration    = errors.New("Dirty database")
	ErrDuplicateData     = errors.New("error duplicate data")
	ErrEncryptData       = errors.New("unable to encrypt data")
	ErrGetData           = errors.New("error get data")
	ErrHashData          = errors.New("unable to hash data")
	ErrInsertData        = errors.New("error insert data")
	ErrInvalidLogin      = errors.New("invalid user")
	ErrInvalidRequest    = errors.New("invalid request")
	ErrInvalidToken      = errors.New("invalid token")
	ErrListenAndServer   = errors.New("unable to listen and serve")
	ErrMigrate           = errors.New("unable to migrate")
	ErrMigrateNoChange   = errors.New("no change")
	ErrMissingData       = errors.New("missing data")
	ErrMissingExpire     = errors.New("invalid expire time")
	ErrParseBody         = errors.New("unable to parse request body JSON")
	ErrParseInt          = errors.New("unable to parse integer")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrReadData          = errors.New("error read data")
	ErrResetPassword     = errors.New("unable to reset password")
	ErrSendEmail         = errors.New("unable to send email")
	ErrTokenExpired      = errors.New("error token expired")
	ErrorConnectDB       = errors.New("unable to connect to database")
)

// HTTP Status
var (
	StatusFatalErr    = "fatal error"
	StatusServerErr   = "server error"
	StatusClientErr   = "client error"
	StatusRedirect    = "redirect"
	StatusSuccess     = "success"
	StatusInformative = "informative"
	StatusUnknown     = "unknown"
)

// App Message
var (
	SuccessGet    = "Success Get"
	SuccessInsert = "Success Insert"
	SuccessUpdate = "Success Update"
	SuccessDelete = "Success Delete"
	SuccessAssign = "Success Assign"
)
