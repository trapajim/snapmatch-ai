package data

import (
	"context"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
)

type ProductData struct {
	appContext  snapmatchai.Context
	productRepo snapmatchai.Repository[*snapmatchai.ProductData]
}

func NewProductData(appContext snapmatchai.Context, productRepo snapmatchai.Repository[*snapmatchai.ProductData]) *ProductData {
	return &ProductData{
		appContext:  appContext,
		productRepo: productRepo,
	}
}

func (c *ProductData) BatchInsert(ctx context.Context, data []snapmatchai.ProductData) error {
	for _, d := range data {
		if err := c.productRepo.Create(ctx, &d); err != nil {
			return err
		}
	}
	return nil
}
