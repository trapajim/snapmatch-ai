package jobworker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/trapajim/snapmatch-ai/server/middleware"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"golang.org/x/sync/errgroup"
	"log"
	"log/slog"
	"path"
	"strings"
	"sync"
	"time"
)

type JobData struct {
	batchPreds *snapmatchai.BatchPrediction
	session    *snapmatchai.Session
}

type ResultHandler interface {
	HandleResult(ctx context.Context, record JSONLRecord) error
}
type JobWorker struct {
	jobChan       chan *JobData
	jobRepo       snapmatchai.Repository[*snapmatchai.BatchPrediction]
	jobBucket     string
	uploader      snapmatchai.Uploader
	handlers      map[string]ResultHandler
	stopChan      chan struct{}
	wg            *sync.WaitGroup
	genAi         snapmatchai.GenAIBatch
	logger        *slog.Logger
	checkInterval time.Duration
}

func NewJobWorker(checkInterval time.Duration, repo snapmatchai.Repository[*snapmatchai.BatchPrediction], logger *slog.Logger, genai snapmatchai.GenAIBatch, uploader snapmatchai.Uploader, jobBucket string) *JobWorker {
	return &JobWorker{
		jobChan:       make(chan *JobData),
		stopChan:      make(chan struct{}),
		jobBucket:     jobBucket,
		handlers:      make(map[string]ResultHandler),
		uploader:      uploader,
		genAi:         genai,
		jobRepo:       repo,
		logger:        logger,
		wg:            &sync.WaitGroup{},
		checkInterval: checkInterval,
	}
}

func (w *JobWorker) RegisterHandler(jobType string, handler ResultHandler) {
	w.handlers[jobType] = handler
}

// Start begins the worker to monitor jobs.
func (w *JobWorker) Start(ctx context.Context) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		activeJobs := make(map[string]*JobData)

		ticker := time.NewTicker(w.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case job := <-w.jobChan:
				activeJobs[job.batchPreds.JobName] = job
				fmt.Printf("Added job: %s\n", job.batchPreds.JobName)
			case <-ticker.C:
				for jobName, job := range activeJobs {
					jobCtx := middleware.SetSession(ctx, *job.session)
					genaiJob, err := w.genAi.GetBatchPredictionJob(jobCtx, job.batchPreds.InternalName)
					if err != nil {
						w.logger.ErrorContext(ctx, fmt.Errorf("error: couldn't get updated job state: %w", err).Error())
						continue
					}
					if genaiJob.Status == job.batchPreds.Status {
						continue
					}
					success := w.handleJobStatusChange(jobCtx, job.batchPreds, genaiJob)
					if success {
						delete(activeJobs, jobName)
					}
				}
			case <-w.stopChan:
				fmt.Println("Stopping worker...")
				return
			}
		}
	}()
}

func (w *JobWorker) handleJobStatusChange(ctx context.Context, job *snapmatchai.BatchPrediction, genaiJob snapmatchai.BatchPrediction) bool {
	job.Status = genaiJob.Status
	job.OutputPath = genaiJob.OutputPath
	err := w.jobRepo.Update(ctx, job)
	if err != nil {
		w.logger.ErrorContext(ctx, fmt.Errorf("error: couldn't update job status: %w", err).Error())
		return false
	}
	if job.Status == "JOB_STATE_PENDING" || job.Status == "JOB_STATE_RUNNING" {
		w.logger.InfoContext(ctx, "Job still running"+job.JobName)
		return false
	}
	if job.Status == "JOB_STATE_FAILED" || job.Status == "JOB_STATE_CANCELLED" {
		w.logger.InfoContext(ctx, "Job failed"+job.JobName)
		return true
	}
	w.logger.InfoContext(ctx, "Job finished"+job.JobName)
	err = w.ReadResults(ctx, job)
	if err != nil {
		w.logger.ErrorContext(ctx, fmt.Errorf("error: couldn't read results: %w", err).Error())
		return true
	}
	return true
}

func (w *JobWorker) ReadResults(ctx context.Context, job *snapmatchai.BatchPrediction) error {
	output := strings.SplitN(strings.TrimPrefix(job.OutputPath, "gs://"), "/", 2)[1]
	f, err := w.uploader.WithBucket(w.jobBucket).GetFile(ctx, path.Join(output, "predictions.jsonl"))
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	handler, ok := w.handlers[job.JobType]
	if !ok {
		return fmt.Errorf("no handler found for job type: %s", job.JobType)
	}
	g, gCtx := errgroup.WithContext(ctx)
	sem := make(chan struct{}, 20)

	for scanner.Scan() {
		line := scanner.Text()
		var record JSONLRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			w.logger.ErrorContext(ctx, fmt.Errorf("failed to parse JSON line: %w", err).Error())
			continue
		}
		sem <- struct{}{}
		g.Go(func() error {
			defer func() { <-sem }() // Release semaphore slot when done
			if err := handler.HandleResult(gCtx, record); err != nil {
				log.Println("error occurred while handling result", err)
			}
			return nil
		})
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

// AddJob sends a new job to the worker for monitoring.
func (w *JobWorker) AddJob(job *snapmatchai.BatchPrediction, session *snapmatchai.Session) {
	w.logger.InfoContext(context.Background(), "Adding job to worker", slog.String("job_name", job.JobName), slog.String("job_id", job.ID))
	w.jobChan <- &JobData{batchPreds: job, session: session}
}

// Stop signals the worker to stop and waits for it to finish.
func (w *JobWorker) Stop() {
	close(w.stopChan)
	w.wg.Wait()
}

func (w *JobWorker) Len() int {
	return len(w.jobChan)
}

type JSONLRecord struct {
	Request struct {
		Contents []struct {
			Parts []struct {
				Text *string `json:"text"`
			} `json:"parts"`
		} `json:"contents"`
	} `json:"request"`
	Response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	} `json:"response"`
}
