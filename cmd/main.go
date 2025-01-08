package main

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/trapajim/snapmatch-ai/server"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
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
	client, err := storage.NewClient(ctx)
	fatalErr(err)
	return &snapmatchai.Context{
		Logger:  slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		Storage: client,
	}
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
