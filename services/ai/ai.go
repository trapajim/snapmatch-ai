package ai

import (
	"context"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
)

type Service struct {
	appContext snapmatchai.Context
}

func NewService(appContext snapmatchai.Context) *Service {
	return &Service{appContext: appContext}
}

func (s *Service) GenerateEmbeddings(ctx context.Context, text string) ([]float32, error) {
	emb, err := s.appContext.GenAI.GenerateEmbeddings(ctx, text)
	if err != nil {
		return nil, err
	}
	return emb, nil
}
