package common

import (
	"log"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}

func NewAppError(err error, message string, code int) *AppError {
	log.Println(err)
	return &AppError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrInvalidInput       = AppError{Code: http.StatusBadRequest, Message: "invalid input"}
	ErrUnauthorized       = AppError{Code: http.StatusUnauthorized, Message: "unauthorized"}
	ErrForbidden          = AppError{Code: http.StatusForbidden, Message: "forbidden"}
	ErrNotFound           = AppError{Code: http.StatusNotFound, Message: "resource not found"}
	ErrInternalServer     = AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	ErrInvalidCredentials = AppError{Code: http.StatusUnauthorized, Message: "invalid credentials"}
	ErrEmailExists        = AppError{Code: http.StatusBadRequest, Message: "email already exists"}
)
