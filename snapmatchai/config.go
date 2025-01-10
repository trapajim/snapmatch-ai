package snapmatchai

import "os"

type Config struct {
	StorageBucket     string
	ProjectID         string
	Location          string
	DatasetID         string
	BQVertexConn      string
	TableID           string
	BQMultiModalModel string
}

func NewConfig() *Config {
	return &Config{
		StorageBucket:     os.Getenv("STORAGE_BUCKET"),
		ProjectID:         os.Getenv("PROJECT_ID"),
		Location:          os.Getenv("LOCATION"),
		DatasetID:         os.Getenv("DATASET_ID"),
		BQVertexConn:      os.Getenv("BQ_VERTEX_CONN"),
		TableID:           os.Getenv("TABLE_ID"),
		BQMultiModalModel: os.Getenv("BQ_MULTI_MODAL_MODEL"),
	}
}
