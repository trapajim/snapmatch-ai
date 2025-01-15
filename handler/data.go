package handler

import (
	"encoding/csv"
	"github.com/trapajim/snapmatch-ai/server"
	"github.com/trapajim/snapmatch-ai/services/data"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"net/http"
)

type DataHandler struct {
	s              *server.Server
	productService *data.ProductData
}

func RegisterDataHandler(s *server.Server, productService *data.ProductData) {
	dataHandler := &DataHandler{s: s, productService: productService}
	s.RegisterRoute("/data", dataHandler.Post)
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
		f := snapmatchai.ProductData{
			Data: record,
		}
		result = append(result, f)
	}
	err = h.productService.BatchInsert(r.Context(), result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
