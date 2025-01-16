package predictions

import (
	"github.com/trapajim/snapmatch-ai/services/ai"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
)

type ProductSearchTerm struct {
	uploader snapmatchai.Uploader
	products []*snapmatchai.ProductData
}

func NewProductSearchTerm(uploader snapmatchai.Uploader, products []*snapmatchai.ProductData) *ProductSearchTerm {
	return &ProductSearchTerm{
		uploader: uploader,
		products: products,
	}
}

func (c *ProductSearchTerm) Name() string {
	return "product_search_term"
}

func (c *ProductSearchTerm) BuildPrediction() []ai.PredictionRequest {
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
	req := make([]ai.PredictionRequest, len(c.products))
	for i, product := range c.products {
		line := c.buildBatchJobLine(product.GetID(), product.Data, instructions)
		req[i] = line
	}
	return req
}

func (c *ProductSearchTerm) buildBatchJobLine(id string, product map[string]string, instructions string) ai.PredictionRequest {
	for key, value := range product {
		instructions += "\n" + key + ": " + value
	}
	line := ai.PredictionRequest{
		Request: ai.Request{
			Contents: []ai.Content{
				{
					Role:  "user",
					Parts: []ai.Parts{{Text: "id::" + id}, {Text: instructions}},
				},
			},
		},
	}

	return line
}
