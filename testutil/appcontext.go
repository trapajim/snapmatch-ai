package snapmatchai

import (
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"github.com/trapajim/snapmatch-ai/testutil/mocks"
	"log/slog"
	"testing"
)

// NewContextForTest creates a new context for testing
// logger defaults to slog.Default()
// storage is a mock uploader
func NewContextForTest(t *testing.T) snapmatchai.Context {
	return snapmatchai.Context{
		Logger:  slog.Default(),
		Storage: mocks.NewMockUploader(t),
		DB:      mocks.NewMockDB(t),
		Config:  snapmatchai.NewConfig(),
	}
}
