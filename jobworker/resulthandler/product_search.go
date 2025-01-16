package resulthandler

import (
	"context"
	"fmt"
	"github.com/trapajim/snapmatch-ai/jobworker"
	"github.com/trapajim/snapmatch-ai/services/asset"
	"github.com/trapajim/snapmatch-ai/services/data"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"log"
	"strings"
)

type ProductSearch struct {
	assetService   *asset.Service
	productService *data.ProductData
}

func NewProductSearch(assetService *asset.Service, prodService *data.ProductData) *ProductSearch {
	return &ProductSearch{
		assetService:   assetService,
		productService: prodService,
	}
}

func (c *ProductSearch) HandleResult(ctx context.Context, record jobworker.JSONLRecord) error {
	var productID string
	for _, part := range record.Request.Contents[0].Parts {
		if part.Text != nil && strings.HasPrefix(*part.Text, "id::") {
			productID = strings.TrimPrefix(*part.Text, "id::")
			break
		}
	}
	var searchTerm string
	log.Println("record.Response.Candidates", productID)
	if len(record.Response.Candidates) > 0 {
		searchTerm = record.Response.Candidates[0].Content.Parts[0].Text
	}
	res, _, err := c.assetService.Search(ctx, strings.TrimSpace(searchTerm), asset.High, snapmatchai.Pagination{})
	if err != nil {
		return fmt.Errorf("failed to search for product: %w", err)
	}
	if len(res) == 0 {
		return fmt.Errorf("no results found for product: %s", productID)
	}
	product, err := c.productService.Get(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}
	product.AssetLinks = []string{res[0].ObjName}
	if err := c.productService.Update(ctx, *product); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}
