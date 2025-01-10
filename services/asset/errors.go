package asset

import (
	"errors"
	"fmt"
	"net/http"
)

var tableCreatedError = errors.New("table created")

type PaginationError struct {
	Err     error
	Message string
	Code    int
}

func NewPaginationError(err error, message string) *PaginationError {
	return &PaginationError{
		Err:     err,
		Message: message,
		Code:    http.StatusBadRequest,
	}
}
func (e *PaginationError) Unwrap() error {
	return e.Err
}
func (e *PaginationError) Error() string {
	return fmt.Sprintf("Pagination Error (code %d): %s", e.Code, e.Message)
}

type UploadError struct {
	Err     error
	Message string
	Code    int
}

func NewUploadError(err error, message string, code int) *UploadError {
	return &UploadError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

func (e *UploadError) Error() string {
	return fmt.Sprintf("Upload Error (code %d): %s", e.Code, e.Message)
}

func (e *UploadError) Unwrap() error {
	return e.Err
}
