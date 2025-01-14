package resulthandler

import (
	"context"
	"fmt"
	"github.com/trapajim/snapmatch-ai/jobworker"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"strings"
)

type ImageCategory struct {
	uploader snapmatchai.Uploader
}

func NewImageCategory(uploader snapmatchai.Uploader) *ImageCategory {
	return &ImageCategory{
		uploader: uploader,
	}
}

func (c *ImageCategory) HandleResult(ctx context.Context, record jobworker.JSONLRecord) error {
	var assetName string
	for _, part := range record.Request.Contents[0].Parts {
		if part.Text != nil && strings.HasPrefix(*part.Text, "file::") {
			assetName = strings.TrimPrefix(*part.Text, "file::")
			assetName = strings.SplitN(strings.TrimPrefix(assetName, "gs://"), "/", 2)[1]
			break
		}
	}
	var cat string
	if len(record.Response.Candidates) > 0 {
		cat = record.Response.Candidates[0].Content.Parts[0].Text
	}
	err := c.uploader.UpdateMetadata(ctx, assetName, map[string]string{"category": strings.TrimSpace(cat)})
	if err != nil {
		return fmt.Errorf("failed to update metadata: %w", err)
	}
	return nil
}
