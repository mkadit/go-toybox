package models

import (
	"errors"
)

type ApplicationErr struct {
	errValue error
	errType  error
}

func NewErr(err, msg error) *ApplicationErr {
	return &ApplicationErr{
		errValue: err,
		errType:  msg,
	}
}

// Error print error message
func (e *ApplicationErr) Error() string {
	return e.errType.Error()
}

// Error print error message
func (e *ApplicationErr) Value() error {
	return e.errValue
}
func (e *ApplicationErr) Is(target error) bool {
	return errors.Is(target, e.errValue)
}

func (e *ApplicationErr) IsType(target error) bool {
	return errors.Is(e.errType, target)
}

// Wrap description
func (e *ApplicationErr) Wrap(err error, msg error) {
	e.errValue = err
	e.errType = msg
}

// Unwrap
func (e *ApplicationErr) Unwrap() error {
	return errors.Join(e.errType, e.errValue)

}
