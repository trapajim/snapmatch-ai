package snapmatchai

import (
	"context"
	"io"
	"log/slog"
	"time"
)

type Uploader interface {
	Upload(ctx context.Context, file io.Reader, object string) error
	WithBucket(bucket string) Uploader
	SignUrl(ctx context.Context, object string, expiry time.Duration) (string, error)
}

type Void struct{}
type DB interface {
	// Query executes a parameterized query and maps the results to a target struct
	// Void should be passed if no result is expected
	Query(ctx context.Context, queryString string, parameters map[string]any, target any) error
	TableExists(ctx context.Context, dataset, table string) error
	Schema(ctx context.Context, dataset, tableName string) ([]DBSchema, error)
}

type GenAIBatch interface {
	CreateBatchPredictionJob(context.Context, BatchPredictionRequest) (BatchPredictionJobConfig, error)
}
type Context struct {
	Logger     *slog.Logger
	Storage    Uploader
	DB         DB
	GenAIBatch GenAIBatch
	Config     *Config
}
