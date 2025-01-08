package main

import (
	"github.com/trapajim/snapmatch-ai/server"
	"log"
	"log/slog"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	s := server.NewServer(":"+port, slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := s.Start(); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
