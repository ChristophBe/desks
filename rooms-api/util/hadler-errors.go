package util

import (
	"fmt"
	"net/http"
)

type HandlerError struct {
	Cause   error  `json:"-"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (h HandlerError) Error() string {
	return fmt.Sprintf("Message: %s, Status: %d, Cauese: %v ", h.Message, h.Status, h.Cause)
}

func InternalServerError(err error) error {
	return HandlerError{
		Cause:   err,
		Message: "unexpected failure",
		Status:  http.StatusInternalServerError,
	}
}

func NotFound(err error) error {
	return HandlerError{
		Cause:   err,
		Message: "not found",
		Status:  http.StatusNotFound,
	}
}
func BadRequest(err error) error {
	return HandlerError{
		Cause:   err,
		Message: "bad request",
		Status:  http.StatusBadRequest,
	}
}
func Unauthorized(err error) error {
	return HandlerError{
		Cause:   err,
		Message: "unauthorized request",
		Status:  http.StatusUnauthorized,
	}
}
