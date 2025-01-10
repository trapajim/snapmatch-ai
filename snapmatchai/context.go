package snapmatchai

import (
	"context"
	"io"
	"log/slog"
)

type Uploader interface {
	Upload(ctx context.Context, file io.Reader, object string) error
}

type Void struct{}
type DB interface {
	// Query executes a parameterized query and maps the results to a target struct
	// Void should be passed if no result is expected
	Query(ctx context.Context, queryString string, parameters map[string]any, target any) error
	TableExists(ctx context.Context, dataset, table string) error
	Schema(ctx context.Context, dataset, tableName string) ([]DBSchema, error)
}

type Context struct {
	Logger  *slog.Logger
	Storage Uploader
	DB      DB
	Config  *Config
}
