package predictions

import (
	"context"
	"encoding/csv"
	"errors"
	"github.com/trapajim/snapmatch-ai/services/ai"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"log"
)

type ProductSearchTerm struct {
	csv      string
	uploader snapmatchai.Uploader
}

func NewProductSearchTerm(csv string, uploader snapmatchai.Uploader) *ProductSearchTerm {
	return &ProductSearchTerm{
		csv:      csv,
		uploader: uploader,
	}
}

func (c *ProductSearchTerm) Name() string {
	return "product_search_term"
}

func (c *ProductSearchTerm) BuildPrediction() []ai.PredictionRequest {
	// read the csv file
	f, err := c.uploader.WithBucket("test").GetFile(context.Background(), c.csv)
	if err != nil {
		errAs := &snapmatchai.Error{}
		if errors.As(err, &errAs) {
			log.Println("Error reading file: ", errAs.Unwrap().Error())
			return nil
		}
		log.Println("Error reading file: ", err)
		return nil
	}
	defer f.Close()
	reader := csv.NewReader(f)
	headers, err := reader.Read()
	if err != nil {
		log.Println("Error reading CSV header: ", err)
		return nil
	}

	rows, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading CSV rows: ", err)
		return nil
	}

	var result []map[string]string
	for _, row := range rows {
		record := make(map[string]string)
		for i, value := range row {
			record[headers[i]] = value
		}
		result = append(result, record)
	}
	instructions := `
		Given the product data below, 
		generate a natural and context-rich search query that summarizes the product effectively for use in a vector-based search system. 
		The query should include the product's category, name, and descriptive attributes, emphasizing its key features and intended use. 
		Avoid technical or overly structured formatting.
		\n
		example:
		Input:
			Category: Watch
			Name: Classic Leather Watch
			Description: Stylish and timeless design.
		Output: A stylish and timeless classic leather watch, perfect for everyday wear.
		Now Generate a search term for the following product:
 		`
	log.Println("Instructions: ", instructions)
	req := make([]ai.PredictionRequest, len(result))
	for i, img := range result {
		line := c.buildBatchJobLine(img, instructions)
		req[i] = line
	}
	return req
}

func (c *ProductSearchTerm) buildBatchJobLine(product map[string]string, instructions string) ai.PredictionRequest {
	for key, value := range product {
		instructions += "\n" + key + ": " + value
	}
	line := ai.PredictionRequest{
		Request: ai.Request{
			Contents: []ai.Content{
				{
					Role:  "user",
					Parts: []ai.Parts{{Text: "name::" + product["Name"]}, {Text: instructions}},
				},
			},
		},
	}

	return line
}
