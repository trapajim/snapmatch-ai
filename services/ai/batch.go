package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/trapajim/snapmatch-ai/jobworker"
	"github.com/trapajim/snapmatch-ai/server/middleware"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"log"
	"log/slog"
	"time"
)

type PredictionBuilder interface {
	Name() string
	BuildPrediction() []PredictionRequest
}
type BatchPredictionService struct {
	appContext snapmatchai.Context
	repo       snapmatchai.Repository[*snapmatchai.BatchPrediction]
	worker     *jobworker.JobWorker
}

func NewBatchPredictionService(appContext snapmatchai.Context, repo snapmatchai.Repository[*snapmatchai.BatchPrediction], worker *jobworker.JobWorker) *BatchPredictionService {
	return &BatchPredictionService{appContext: appContext, repo: repo, worker: worker}
}

func (b *BatchPredictionService) Predict(ctx context.Context, builder PredictionBuilder) error {
	session := middleware.GetSession(ctx)
	if session == nil {
		return snapmatchai.NewError(errors.New("session not found"), "session not found", 404)
	}
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
	input := fmt.Sprintf("gs://%s/%s/%s", b.appContext.Config.JobsStorageBucket, session.SessionID(), jobName+".json")
	output := fmt.Sprintf("gs://%s/%s/result", b.appContext.Config.JobsStorageBucket, session.SessionID())
	log.Println("input", input)
	log.Println("output", output)
	job, err := b.createBatchPredictionJob(ctx, jobName, input, output)
	if err != nil {
		return snapmatchai.NewError(err, "failed to create batch prediction job", 500)
	}
	job.JobType = builder.Name()
	err = b.repo.Create(ctx, &job)
	if err != nil {
		log.Println("failed to save batch prediction job", err)
		return snapmatchai.NewError(err, "failed to save batch prediction job", 500)
	}
	b.worker.AddJob(&job, &session)
	b.appContext.Logger.InfoContext(ctx, "Batch prediction job created", slog.String("job_id", job.ID), slog.String("job_name", job.JobName))
	return nil
}

func (b *BatchPredictionService) createBatchPredictionJob(ctx context.Context, name, inputPath, outputPath string) (snapmatchai.BatchPrediction, error) {
	modelName := "gemini-1.5-flash-002"
	modelParameters := map[string]any{
		"temperature": 0.2,
	}
	request := snapmatchai.NewBatchPrediction(name, modelName, inputPath, outputPath, modelParameters)
	job, err := b.appContext.GenAIBatch.CreateBatchPredictionJob(ctx, *request)
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
		return snapmatchai.BatchPrediction{}, err
	}
	return job, nil
}
