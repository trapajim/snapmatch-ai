package handler

import (
	"github.com/trapajim/snapmatch-ai/server"
	"github.com/trapajim/snapmatch-ai/services/job"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"github.com/trapajim/snapmatch-ai/templates/models"
	"github.com/trapajim/snapmatch-ai/templates/pages"
	"net/http"
)

type JobsHandler struct {
	s          *server.Server
	jobService *job.Job
}

func RegisterJobsHandler(s *server.Server, jobService *job.Job) {
	jobsHandler := &JobsHandler{s: s, jobService: jobService}
	s.RegisterRoute("/jobs", jobsHandler.Get)
}
func (h *JobsHandler) Get(w http.ResponseWriter, r *http.Request) {
	jobsList, err := h.jobService.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = pages.Jobs(jobsToViewModel(jobsList)).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func jobsToViewModel(jobs []*snapmatchai.BatchPrediction) []models.Job {
	templateJobs := make([]models.Job, len(jobs))

	for i, j := range jobs {

		templateJobs[i] = models.Job{
			ID:   j.ID,
			Name: j.JobName,
		}
		switch j.Status {
		case "JOB_STATE_RUNNING":
			templateJobs[i].Status = "Running"
		case "JOB_STATE_PENDING":
			templateJobs[i].Status = "Pending"
		case "JOB_STATE_FAILED":
			templateJobs[i].Status = "Failed"
		case "JOB_STATE_SUCCEEDED":
			templateJobs[i].Status = "Success"
		case "JOB_STATE_CANCELLED":
			templateJobs[i].Status = "Failed"
		}
	}
	return templateJobs
}
