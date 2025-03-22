package errors

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeAuth         ErrorType = "AUTHENTICATION_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
	ErrorTypeBadRequest   ErrorType = "BAD_REQUEST"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeConflict     ErrorType = "CONFLICT_ERROR"
)

type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Code    int       `json:"code"`
	Errors  any       `json:"errors,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

var statusCodeMap = map[ErrorType]int{
	ErrorTypeValidation:   http.StatusBadRequest,
	ErrorTypeAuth:         http.StatusUnauthorized,
	ErrorTypeNotFound:     http.StatusNotFound,
	ErrorTypeInternal:     http.StatusInternalServerError,
	ErrorTypeBadRequest:   http.StatusBadRequest,
	ErrorTypeUnauthorized: http.StatusUnauthorized,
}

func ValidationError(message string, errors any) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Code:    statusCodeMap[ErrorTypeValidation],
		Errors:  errors,
	}
}

func AuthError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeAuth,
		Message: message,
		Code:    statusCodeMap[ErrorTypeAuth],
	}
}

func NotFoundError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Code:    statusCodeMap[ErrorTypeNotFound],
	}
}

func InternalError(err error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: "Internal server error",
		Code:    statusCodeMap[ErrorTypeInternal],
		Errors:  fmt.Sprintf("%v", err),
	}
}

func BadRequestError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
		Code:    statusCodeMap[ErrorTypeBadRequest],
	}
}

func ConflictError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
		Code:    http.StatusConflict,
	}
}
