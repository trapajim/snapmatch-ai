package snapmatchai

import "os"

type Config struct {
	StorageBucket string
}

func NewConfig() *Config {
	return &Config{
		StorageBucket: os.Getenv("STORAGE_BUCKET"),
	}
}
