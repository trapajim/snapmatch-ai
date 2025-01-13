package snapmatchai

import (
	"cloud.google.com/go/firestore"
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
type DataStore interface {
	// Query executes a parameterized query and maps the results to a target struct
	// Void should be passed if no result is expected
	Query(ctx context.Context, queryString string, parameters map[string]any, target any) error
	TableExists(ctx context.Context, dataset, table string) error
	Schema(ctx context.Context, dataset, tableName string) ([]DBSchema, error)
}
type Entity interface {
	GetID() string
	SetID(id string)
}
type Repository[T Entity] interface {
	Create(ctx context.Context, entity T) error
	Read(ctx context.Context, id string) (T, error)
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters map[string]any) ([]T, error)
}

type GenAIBatch interface {
	CreateBatchPredictionJob(context.Context, BatchPrediction) (BatchPrediction, error)
}
type Context struct {
	Logger     *slog.Logger
	Storage    Uploader
	FireStore  *firestore.Client
	DB         DataStore
	GenAIBatch GenAIBatch
	Config     *Config
}
