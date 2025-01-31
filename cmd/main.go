package main

import (
	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	googleGenAI "github.com/google/generative-ai-go/genai"
	"github.com/trapajim/snapmatch-ai/datastore"
	"github.com/trapajim/snapmatch-ai/genai"
	"github.com/trapajim/snapmatch-ai/handler"
	"github.com/trapajim/snapmatch-ai/jobworker"
	"github.com/trapajim/snapmatch-ai/jobworker/resulthandler"
	_ "github.com/trapajim/snapmatch-ai/memory"
	"github.com/trapajim/snapmatch-ai/server"
	"github.com/trapajim/snapmatch-ai/server/middleware"
	"github.com/trapajim/snapmatch-ai/services/ai"
	"github.com/trapajim/snapmatch-ai/services/asset"
	"github.com/trapajim/snapmatch-ai/services/data"
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
	aiService := ai.NewService(appContext)
	assetService := asset.NewService(appContext)
	productService := data.NewProductData(appContext, datastore.NewFirestoreRepository[*snapmatchai.ProductData](appContext.FireStore, aiService, "product_data"))
	batchPredictionRepository := datastore.NewFirestoreRepository[*snapmatchai.BatchPrediction](appContext.FireStore, aiService, "batch_predictions")
	worker := jobworker.NewJobWorker(20*time.Second, batchPredictionRepository, appContext.Logger, appContext.GenAIBatch, appContext.Storage, appContext.Config.JobsStorageBucket)
	worker.RegisterHandler("categorize_images", resulthandler.NewImageCategory(appContext.Storage))
	worker.RegisterHandler("product_search_term", resulthandler.NewProductSearch(assetService, productService))
	worker.Start(context.Background())
	batchPredictionService := ai.NewBatchPredictionService(appContext, batchPredictionRepository, worker)
	jobService := job.NewService(appContext, batchPredictionRepository)
	s.RegisterMiddleware(middleware.AuthMiddleware(appContext.SessionManager))
	handler.RegisterIndexHandler(s, jobService)
	handler.RegisterAssetsHandler(s, assetService, batchPredictionService, appContext.Storage)
	handler.RegisterJobsHandler(s, jobService)
	handler.RegisterDataHandler(s, productService, aiService, batchPredictionService, appContext.Storage)
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
	googleGenaiClient, err := googleGenAI.NewClient(context.Background(), option.WithAPIKey(config.GeminiAPIKey))
	fatalErr(err)
	genaiClient := genai.NewClient(googleGenaiClient)
	sessionManager, err := snapmatchai.NewManager("memory", "gosessionid", 3600)
	fatalErr(err)
	go sessionManager.GC()
	return snapmatchai.Context{
		Logger:         slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		Storage:        client,
		FireStore:      firestoreClient,
		DB:             datastore.NewBigQuery(bqClient),
		GenAIBatch:     aiBatchClient,
		GenAI:          genaiClient,
		Config:         config,
		SessionManager: sessionManager,
	}
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
