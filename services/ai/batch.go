package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"log/slog"
	"time"
)

type PredictionBuilder interface {
	Name() string
	BuildPrediction() []PredictionRequest
}
type BatchPredictionService struct {
	appContext snapmatchai.Context
}

func NewBatchPredictionService(appContext snapmatchai.Context) *BatchPredictionService {
	return &BatchPredictionService{appContext: appContext}
}

func (b *BatchPredictionService) Predict(ctx context.Context, builder PredictionBuilder) error {
	reqs := builder.BuildPrediction()
	var buffer bytes.Buffer
	for _, req := range reqs {
		l, err := json.Marshal(req)
		if err != nil {
			b.appContext.Logger.ErrorContext(ctx, "error occurred while marshalling request", slog.String("error", err.Error()))
			continue // Skip this request and move to the next.
		}
		buffer.Write(l)
		buffer.Write([]byte("\n"))
	}
	jobName := fmt.Sprintf("jobs-%d", time.Now().UTC().Unix())
	err := b.appContext.Storage.WithBucket(b.appContext.Config.JobsStorageBucket).Upload(ctx, &buffer, jobName+".json")
	if err != nil {
		return fmt.Errorf("failed to write batch job: %w", err)
	}
	input := fmt.Sprintf("gs://%s/%s", b.appContext.Config.JobsStorageBucket, jobName+".json")
	output := fmt.Sprintf("gs://%s/result.json", b.appContext.Config.JobsStorageBucket)
	job, err := b.createBatchPredictionJob(ctx, jobName, input, output)
	if err != nil {
		return fmt.Errorf("failed to create batch prediction job: %w", err)
	}
	b.appContext.Logger.InfoContext(ctx, "Batch prediction job created", slog.String("job_id", job.Name))
	return nil
}

func (b *BatchPredictionService) createBatchPredictionJob(ctx context.Context, name, inputPath, outputPath string) (snapmatchai.BatchPredictionJobConfig, error) {
	modelName := "gemini-1.5-flash-002"
	modelParameters := map[string]any{
		"temperature": 0.2,
	}
	request := snapmatchai.NewBatchPredictionRequest(name, modelName, inputPath, outputPath, modelParameters)
	job, err := b.appContext.GenAIBatch.CreateBatchPredictionJob(ctx, request)
	if err != nil {
		errAs := &snapmatchai.Error{}
		if errors.As(err, &errAs) {
			b.appContext.Logger.ErrorContext(ctx, "Service: Could not create batch prediction job",
				slog.Int("status_code", errAs.Code),
				slog.String("error", errAs.Error()),
				slog.String("message", errAs.Message),
				slog.String("unwrapped error", errAs.Unwrap().Error()),
			)
		} else {
			b.appContext.Logger.ErrorContext(ctx, "unable to create batch prediction job", slog.String("error", err.Error()))
		}
		return snapmatchai.BatchPredictionJobConfig{}, err
	}
	return job, nil
}
