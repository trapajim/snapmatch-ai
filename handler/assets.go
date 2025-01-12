package handler

import (
	"fmt"
	"github.com/trapajim/snapmatch-ai/server"
	"github.com/trapajim/snapmatch-ai/services/ai"
	"github.com/trapajim/snapmatch-ai/services/ai/predictions"
	"github.com/trapajim/snapmatch-ai/services/asset"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"github.com/trapajim/snapmatch-ai/templates/models"
	"github.com/trapajim/snapmatch-ai/templates/pages"
	"github.com/trapajim/snapmatch-ai/templates/partials"
	"log"
	"net/http"
)

type AssetsHandler struct {
	s            *server.Server
	service      *asset.Service
	batchService *ai.BatchPredictionService
}

func RegisterAssetsHandler(server *server.Server, service *asset.Service, aiService *ai.BatchPredictionService) {
	handler := &AssetsHandler{service: service, batchService: aiService}
	server.RegisterRoute("GET /assets", handler.Get)
	server.RegisterRoute("POST /assets", handler.Upload)
	server.RegisterRoute("POST /asset/predict", handler.Predict)

}

func (h *AssetsHandler) Predict(w http.ResponseWriter, r *http.Request) {
	assets, _, err := h.service.List(r.Context(), snapmatchai.Pagination{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.batchService.Predict(r.Context(), predictions.NewCategorizeImages(assets))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AssetsHandler) Get(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	var assets []snapmatchai.FileRecord
	var err error
	log.Println(q)
	if q != "" {
		assets, _, err = h.service.Search(r.Context(), q, snapmatchai.Pagination{})
	} else {
		assets, _, err = h.service.List(r.Context(), snapmatchai.Pagination{})
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	viewModel := make(models.Assets, len(assets))
	for i, a := range assets {
		viewModel[i] = models.Asset{
			Name: a.ObjName,
			Size: a.Size,
			Type: a.ContentType,
			Date: a.Updated,
			URI:  a.SignedURL,
		}
	}
	if r.Header.Get("HX-Request") == "true" {
		err = partials.Assets(viewModel).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	err = pages.Assets(viewModel).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AssetsHandler) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1000 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	fileChan := make(chan asset.BatchUploadRequest)
	errChan := make(chan error, 1)
	go func() {
		defer close(fileChan)
		for _, file := range files {
			f, err := file.Open()
			if err != nil {
				errChan <- fmt.Errorf("failed to open file: %v", err)
				return
			}
			fileChan <- asset.BatchUploadRequest{
				File: f,
				Name: file.Filename,
			}
		}
	}()
	go func() {
		err := h.service.BatchUpload(r.Context(), fileChan)
		if err != nil {
			errChan <- fmt.Errorf("batch upload failed: %v", err)
		} else {
			errChan <- nil
		}
	}()
	select {
	case err := <-errChan:
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Info("Successfully uploaded files").SetHXTriggerHeader(w)
		return
	case <-r.Context().Done():
		http.Error(w, "Request was canceled", http.StatusRequestTimeout)
		return
	}
}
