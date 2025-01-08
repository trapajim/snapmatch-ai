package snapmatchai

import (
	"cloud.google.com/go/storage"
	"log/slog"
)

type Config struct {
	StorageBucket string
}
type Context struct {
	Logger  *slog.Logger
	Storage *storage.Client
	Config  *Config
}
