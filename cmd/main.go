package main

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/trapajim/snapmatch-ai/server"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"github.com/trapajim/snapmatch-ai/uploader"
	"log"
	"log/slog"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	appContext := createContext(context.Background())
	s := server.NewServer(":"+port, appContext.Logger)
	if err := s.Start(); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}

func createContext(ctx context.Context) *snapmatchai.Context {
	config := snapmatchai.NewConfig()
	storageClient, err := storage.NewGRPCClient(ctx)
	fatalErr(err)
	client := uploader.NewUploader(storageClient, config.StorageBucket)
	fatalErr(err)

	return &snapmatchai.Context{
		Logger:  slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		Storage: client,
		Config:  config,
	}
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
