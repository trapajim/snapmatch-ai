package handler

import (
	"github.com/trapajim/snapmatch-ai/server"
	"github.com/trapajim/snapmatch-ai/templates/pages"
	"net/http"
)

type IndexHandler struct {
	s *server.Server
}

func RegisterIndexHandler(s *server.Server) {
	idxHandler := &IndexHandler{s: s}
	s.RegisterRoute("/", idxHandler.Get)
}
func (h *IndexHandler) Get(w http.ResponseWriter, r *http.Request) {
	err := pages.Dashboard("title").Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
