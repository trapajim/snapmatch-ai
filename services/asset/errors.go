package asset

import "fmt"

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
