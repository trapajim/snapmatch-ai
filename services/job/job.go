package job

import (
	"context"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
)

type Job struct {
	appContext snapmatchai.Context
	jobRepo    snapmatchai.Repository[*snapmatchai.BatchPrediction]
}

func NewService(appContext snapmatchai.Context, jobRepo snapmatchai.Repository[*snapmatchai.BatchPrediction]) *Job {
	return &Job{appContext: appContext, jobRepo: jobRepo}
}

type JobStats struct {
	Total   int
	Running int
	Failed  int
	Success int
}

func (j *Job) List(ctx context.Context) ([]*snapmatchai.BatchPrediction, error) {
	return j.jobRepo.List(ctx, nil)
}
func (j *Job) Stats(ctx context.Context) (JobStats, error) {
	jobs, err := j.jobRepo.List(ctx, nil)
	if err != nil {
		return JobStats{}, err
	}

	stats := JobStats{}
	for _, job := range jobs {
		switch job.Status {
		case "JOB_STATE_RUNNING":
			stats.Running++
		case "JOB_STATE_PENDING":
			stats.Running++
		case "JOB_STATE_FAILED":
			stats.Failed++
		case "JOB_STATE_SUCCEEDED":
			stats.Success++
		case "JOB_STATE_CANCELLED":
			stats.Failed++
		}
		stats.Total++
	}

	return stats, nil

}
