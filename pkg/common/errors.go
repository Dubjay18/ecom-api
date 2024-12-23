package common

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}

var (
	ErrInvalidInput   = AppError{Code: http.StatusBadRequest, Message: "invalid input"}
	ErrUnauthorized   = AppError{Code: http.StatusUnauthorized, Message: "unauthorized"}
	ErrForbidden      = AppError{Code: http.StatusForbidden, Message: "forbidden"}
	ErrNotFound       = AppError{Code: http.StatusNotFound, Message: "resource not found"}
	ErrInternalServer = AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
)
