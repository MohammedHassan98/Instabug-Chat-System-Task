package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   error  `json:"-"` // Internal error details (not exposed to client)
}

func (e *AppError) String() string {
	if e.Error != nil {
		return fmt.Sprintf("Error: %s (Code: %d, Internal: %v)", e.Message, e.Code, e.Error)
	}
	return fmt.Sprintf("Error: %s (Code: %d)", e.Message, e.Code)
}

// Common errors
var (
	ErrInvalidInput = func(err error) *AppError {
		return &AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid input provided",
			Error:   err,
		}
	}

	ErrNotFound = func(resource string) *AppError {
		return &AppError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("%s not found", resource),
		}
	}

	ErrUnauthorized = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized access",
	}

	ErrInternalServer = func(err error) *AppError {
		return &AppError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Error:   err,
		}
	}
)
