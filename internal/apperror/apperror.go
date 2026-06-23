package apperror

import "net/http"

type AppError struct {
	Code    int
	Status  string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func BadRequest(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Status:  "VALIDATION_ERROR",
		Message: message,
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Status:  "NOT_FOUND",
		Message: message,
	}
}

func InternalServerError(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Status:  "INTERNAL_SERVER_ERROR",
		Message: message,
	}
}
