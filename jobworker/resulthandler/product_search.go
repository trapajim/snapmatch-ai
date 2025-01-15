package resulthandler

import (
	"context"
	"fmt"
	"github.com/trapajim/snapmatch-ai/jobworker"
	"github.com/trapajim/snapmatch-ai/services/asset"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"log"
	"strings"
)

type ProductSearch struct {
	assetService *asset.Service
}

func NewProductSearch(assetService *asset.Service) *ProductSearch {
	return &ProductSearch{
		assetService: assetService,
	}
}

func (c *ProductSearch) HandleResult(ctx context.Context, record jobworker.JSONLRecord) error {
	var productName string
	for _, part := range record.Request.Contents[0].Parts {
		if part.Text != nil && strings.HasPrefix(*part.Text, "name::") {
			productName = strings.TrimPrefix(*part.Text, "name::")
			break
		}
	}

	var searchTerm string
	if len(record.Response.Candidates) > 0 {
		searchTerm = record.Response.Candidates[0].Content.Parts[0].Text
	}
	res, _, err := c.assetService.Search(ctx, strings.TrimSpace(searchTerm), asset.High, snapmatchai.Pagination{})
	if err != nil {
		return fmt.Errorf("failed to search for product: %w", err)
	}
	if len(res) > 0 {
		log.Println(productName)
		log.Println("Product found: ", res[0].ObjName, " with search term: ", searchTerm)
	}
	return nil
}
