package snapmatchai

import (
	"context"
	"log/slog"
)

type Error struct {
	Err     error
	Message string
	Code    int
}

func NewError(err error, message string, code int) *Error {
	return &Error{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Log(ctx context.Context, logger slog.Logger) {
	logger.ErrorContext(ctx, "Google API error occurred, during file upload",
		slog.Int("status_code", e.Code),
		slog.String("error_body", e.Message),
		slog.String("error", e.Error()),
		slog.String("message", e.Message),
		slog.Any("error", e.Unwrap()),
	)
}
