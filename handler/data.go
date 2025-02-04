package handler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/trapajim/snapmatch-ai/server"
	"github.com/trapajim/snapmatch-ai/services/ai"
	"github.com/trapajim/snapmatch-ai/services/ai/predictions"
	"github.com/trapajim/snapmatch-ai/services/data"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"github.com/trapajim/snapmatch-ai/templates/pages"
	"log"
	"net/http"
	"time"
)

type DataHandler struct {
	s              *server.Server
	productService *data.ProductData
	aiService      *ai.Service
	batchService   *ai.BatchPredictionService
	uploader       snapmatchai.Uploader
}

func RegisterDataHandler(s *server.Server, productService *data.ProductData, aiService *ai.Service, batchService *ai.BatchPredictionService, uploader snapmatchai.Uploader) {
	dataHandler := &DataHandler{s: s, productService: productService, batchService: batchService, aiService: aiService, uploader: uploader}
	s.RegisterRoute("POST /data", dataHandler.Post)
	s.RegisterRoute("GET /data", dataHandler.List)
	s.RegisterRoute("POST /data/match", dataHandler.Match)
}

func (h *DataHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	var products []*snapmatchai.ProductData
	var err error
	if q != "" {
		products, err = h.productService.Search(r.Context(), q)
	} else {
		products, err = h.productService.List(r.Context())
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	productsMap := make([]map[string]string, len(products))
	headers := []string{"Name", "Description", "__asset_url"}
	for i, p := range products {
		imgUrl := ""
		if len(p.AssetLinks) != 0 {
			img, err := h.uploader.SignUrl(r.Context(), p.AssetLinks[0], 10*time.Minute)
			if err != nil {
				log.Println(err)
			} else {
				imgUrl = img
			}
		}
		p.Data["__asset_url"] = imgUrl
		productsMap[i] = p.Data
	}

	if err := pages.Products(productsMap, headers, GetSessionExpiry(r.Context())).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *DataHandler) Match(w http.ResponseWriter, r *http.Request) {
	assets, err := h.productService.List(r.Context())
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	err = h.batchService.Predict(r.Context(), predictions.NewProductSearchTerm(h.uploader, assets))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info("Successfully started prediction job").SetHXTriggerHeader(w)
}

func (h *DataHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1000 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["file"]
	if len(files) != 1 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}
	// verify it's a csv file
	if files[0].Header.Get("Content-Type") != "text/csv" {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}
	file, err := files[0].Open()
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		http.Error(w, "Failed to read CSV header", http.StatusInternalServerError)
		return
	}
	err = validateHeaders(headers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rows, err := reader.ReadAll()
	if err != nil {
		http.Error(w, "Failed to read CSV rows", http.StatusInternalServerError)
		return
	}
	var result []snapmatchai.ProductData
	for _, row := range rows {
		record := make(map[string]string)
		for i, value := range row {
			record[headers[i]] = value
		}
		emb, err := h.aiService.GenerateEmbeddings(r.Context(), record["Name"]+record["Description"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		f := snapmatchai.ProductData{
			Data:       record,
			VectorData: emb,
		}
		result = append(result, f)
	}
	err = h.productService.BatchInsert(r.Context(), result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

const (
	FlagName        int8 = 1 << 0 // 00000001
	FlagDescription int8 = 1 << 1 // 00000010
)

func validateHeaders(headers []string) error {
	var flags int8
	for _, h := range headers {
		if h == "Name" {
			flags |= FlagName
		} else if h == "Description" {
			flags |= FlagDescription
		}
	}

	missing := make([]string, 0)
	if flags&FlagName == 0 {
		missing = append(missing, "Name")
	}
	if flags&FlagDescription == 0 {
		missing = append(missing, "Description")
	}

	if len(missing) > 0 {
		return errors.New("missing headers: " + fmt.Sprint(missing))
	}
	return nil
}
