package main

import (
	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/trapajim/snapmatch-ai/datastore"
	"github.com/trapajim/snapmatch-ai/genai"
	"github.com/trapajim/snapmatch-ai/handler"
	"github.com/trapajim/snapmatch-ai/jobworker"
	"github.com/trapajim/snapmatch-ai/server"
	"github.com/trapajim/snapmatch-ai/services/ai"
	"github.com/trapajim/snapmatch-ai/services/asset"
	"github.com/trapajim/snapmatch-ai/services/job"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"github.com/trapajim/snapmatch-ai/uploader"
	"time"

	"google.golang.org/api/option"
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
	asserService := asset.NewService(appContext)
	batchPredictionRepository := datastore.NewFirestoreRepository[*snapmatchai.BatchPrediction](appContext.FireStore, "batch_predictions")
	worker := jobworker.NewJobWorker(20*time.Second, batchPredictionRepository, appContext.Logger, appContext.GenAIBatch, appContext.Storage, appContext.Config.JobsStorageBucket)
	worker.Start(context.Background())
	batchPredictionService := ai.NewBatchPredictionService(appContext, batchPredictionRepository, worker)
	jobService := job.NewService(appContext, batchPredictionRepository)
	handler.RegisterIndexHandler(s, jobService)
	handler.RegisterAssetsHandler(s, asserService, batchPredictionService)
	if err := s.Start(); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}

func createContext(ctx context.Context) snapmatchai.Context {
	config := snapmatchai.NewConfig()
	storageClient, err := storage.NewClient(ctx)
	fatalErr(err)
	client := uploader.NewUploader(storageClient, config.StorageBucket)
	fatalErr(err)
	bqClient, err := bigquery.NewClient(ctx, config.ProjectID)
	fatalErr(err)
	apiEndpoint := fmt.Sprintf("asia-northeast1-aiplatform.googleapis.com:443")
	aiClient, err := aiplatform.NewJobClient(context.Background(), option.WithEndpoint(apiEndpoint))
	fatalErr(err)
	aiBatchClient := genai.NewBatchClient(aiClient, config.Location, config.ProjectID)
	firestoreClient, err := firestore.NewClient(ctx, config.ProjectID)
	fatalErr(err)
	return snapmatchai.Context{
		Logger:     slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		Storage:    client,
		FireStore:  firestoreClient,
		DB:         datastore.NewBigQuery(bqClient),
		GenAIBatch: aiBatchClient,
		Config:     config,
	}
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
