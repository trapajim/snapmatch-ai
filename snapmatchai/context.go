package snapmatchai

import (
	"context"
	"github.com/trapajim/snapmatch-ai/snapmatchai/mocks"
	"io"
	"log/slog"
	"testing"
)

type Uploader interface {
	Upload(ctx context.Context, file io.Reader, object string) error
}

type Context struct {
	Logger  *slog.Logger
	Storage Uploader
	Config  *Config
}

// NewContextForTest creates a new context for testing
// logger defaults to slog.Default()
// storage is a mock uploader
func NewContextForTest(t *testing.T) Context {
	return Context{
		Logger:  slog.Default(),
		Storage: mocks.NewMockUploader(t),
		Config:  NewConfig(),
	}
}
