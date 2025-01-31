package handler

import (
	"github.com/trapajim/snapmatch-ai/server"
	"github.com/trapajim/snapmatch-ai/server/middleware"
	"github.com/trapajim/snapmatch-ai/services/job"
	"github.com/trapajim/snapmatch-ai/templates/models"
	"github.com/trapajim/snapmatch-ai/templates/pages"
	"log"
	"net/http"
)

type IndexHandler struct {
	s          *server.Server
	jobService *job.Job
}

func RegisterIndexHandler(s *server.Server, jobService *job.Job) {
	idxHandler := &IndexHandler{s: s, jobService: jobService}
	s.RegisterRoute("/", idxHandler.Get)
}
func (h *IndexHandler) Get(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSession(r)
	if session != nil {
		log.Println(session.SessionID())
	} else {
		log.Println("No session found")
	}
	stats, err := h.jobService.Stats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := models.BarChart{
		Labels: []string{"Total", "Success", "Running", "Failed"},
		Data:   []int{stats.Total, stats.Success, stats.Running, stats.Failed},
		Colors: []string{"#4CAF50", "#2196F3", "#FFC107", "#F44336"},
	}
	err = pages.Dashboard("title", data).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
