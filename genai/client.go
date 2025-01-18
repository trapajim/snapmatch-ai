package genai

import (
	"context"
	googleGenAI "github.com/google/generative-ai-go/genai"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
)

type Client struct {
	client *googleGenAI.Client
}

func NewClient(c *googleGenAI.Client) *Client {
	return &Client{client: c}
}

func (c *Client) GenerateEmbeddings(ctx context.Context, text string) ([]float32, error) {
	em := c.client.EmbeddingModel("text-embedding-004")
	emb, err := em.EmbedContent(ctx, googleGenAI.Text(text))
	if err != nil {
		return nil, snapmatchai.NewError(err, "failed to generate embeddings", 400)
	}
	return emb.Embedding.Values, nil
}
