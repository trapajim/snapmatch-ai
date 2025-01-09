package asset

import (
	"context"
	"errors"
	"fmt"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"io"
	"log/slog"
)

type Service struct {
	appContext snapmatchai.Context
}

// NewService creates a new instance of the asset service
// the asset service is responsible for uploading files to storage
func NewService(appContext snapmatchai.Context) *Service {
	return &Service{
		appContext: appContext,
	}
}

// Upload uploads a file to storage
func (s *Service) Upload(ctx context.Context, file io.Reader, fileName string) error {
	err := s.appContext.Storage.Upload(ctx, file, fileName)
	if err != nil {
		var e *snapmatchai.Error
		if ok := errors.As(err, &e); ok {
			s.appContext.Logger.ErrorContext(ctx, "Google API error occurred, during file upload",
				slog.Int("status_code", e.Code),
				slog.String("error", e.Error()),
				slog.String("message", e.Message),
			)
			return NewUploadError(e, e.Error(), e.Code)
		}
		return snapmatchai.NewError(err, fmt.Errorf("could not upload file %s with error %w", fileName, err).Error(), 500)
	}
	return nil
}
