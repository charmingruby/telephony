package validation

import (
	"fmt"
)

func NewInternalErr() error {
	return &ErrInternal{
		Message: "internal error",
	}
}

type ErrInternal struct {
	Message string `json:"message"`
}

func (e *ErrInternal) Error() string {
	return e.Message
}

func NewNotFoundErr(entity string) error {
	return &ErrNotFound{
		Message: fmt.Sprintf("%s not found", entity),
	}
}

type ErrNotFound struct {
	Message string `json:"message"`
}

func (e *ErrNotFound) Error() string {
	return e.Message
}

func NewBadRequestErr(msg string) error {
	return &ErrBadRequest{
		Message: msg,
	}
}

type ErrBadRequest struct {
	Message string `json:"message"`
}

func (e *ErrBadRequest) Error() string {
	return e.Message
}

func NewInvalidCredentialsErr() error {
	return &ErrInvalidCredentials{
		Message: "invalid credentials",
	}
}

type ErrInvalidCredentials struct {
	Message string `json:"message"`
}

func (e *ErrInvalidCredentials) Error() string {
	return e.Message
}

func NewUnauthorizedErr() error {
	return &ErrUnathorized{
		Message: "user don't have necessary permissions to do this action",
	}
}

type ErrUnathorized struct {
	Message string `json:"message"`
}

func (e *ErrUnathorized) Error() string {
	return e.Message
}

func NewConflictErr(entity, field string) error {
	return &ErrConflict{
		Message: fmt.Sprintf("%s %s already taken", entity, field),
	}
}

type ErrConflict struct {
	Message string `json:"message"`
}

func (e *ErrConflict) Error() string {
	return e.Message
}
