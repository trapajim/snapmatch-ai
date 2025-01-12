package predictions

import (
	"github.com/trapajim/snapmatch-ai/services/ai"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
)

type CategorizeImages struct {
	imgs []snapmatchai.FileRecord
}

func NewCategorizeImages(imgs []snapmatchai.FileRecord) *CategorizeImages {
	return &CategorizeImages{imgs: imgs}
}

func (c *CategorizeImages) Name() string {
	return "categorize_images"
}

func (c *CategorizeImages) BuildPrediction() []ai.PredictionRequest {
	req := make([]ai.PredictionRequest, len(c.imgs))
	instructions := `
		Categorize the product image based on its primary content
		1. Analyze the main object(s) in the image.
		2. Assign the most appropriate category based on the dominant content
		3. If the image contains mixed categories (e.g., food on a table), prioritize the primary focus of the image
		4. If unclear or unidentifiable, assign the category Other.
	Guidelines \n
		1. If the image includes multiple categories, select the category that occupies the largest area or appears most prominently.
		2. Avoid assumptions based on context or metadata; rely solely on the visual elements.
	\n
		Return only the primary category name as a single word e.g. food 
 `
	for i, img := range c.imgs {
		line := c.buildBatchJobLine(img, instructions)
		req[i] = line
	}
	return req
}

func (c *CategorizeImages) buildBatchJobLine(product snapmatchai.FileRecord, instructions string) ai.PredictionRequest {
	line := ai.PredictionRequest{
		Request: ai.Request{
			Contents: []ai.Content{
				{
					Role:  "user",
					Parts: []ai.Parts{{Text: "file::" + product.URI}, {Text: instructions}, {FileData: &ai.FileData{FileUri: product.URI, MIMEType: product.ContentType}}},
				},
			},
		},
	}

	return line
}
