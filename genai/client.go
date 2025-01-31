package genai

import (
	"context"
	googleGenAI "github.com/google/generative-ai-go/genai"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"log"
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

func (c *Client) StartChat(ctx context.Context, funcs []snapmatchai.AITools) snapmatchai.AIChat {
	model := c.client.GenerativeModel("gemini-1.5-flash")
	if len(funcs) > 0 {
		model.Tools = []*googleGenAI.Tool{c.aiToolsToGenAI(funcs)}
	}
	chatSession := model.StartChat()
	aiChat := &AIChat{session: chatSession, model: model}
	return aiChat
}

type AIChat struct {
	session *googleGenAI.ChatSession
	model   *googleGenAI.GenerativeModel
}

func (a *AIChat) SendMessage(ctx context.Context, p ...snapmatchai.AIPart) ([]snapmatchai.AIResponse, error) {
	var parts []googleGenAI.Part
	for _, part := range p {
		switch p := part.(type) {
		case snapmatchai.Text:
			parts = append(parts, googleGenAI.Text(p))
		case snapmatchai.Blob:
			parts = append(parts, googleGenAI.Blob{
				MIMEType: p.MIMEType,
				Data:     p.Data,
			})
		case snapmatchai.FunctionCallResponse:
			args := make(map[string]interface{}, len(p.Args))
			for _, arg := range p.Args {
				args[arg.Name] = arg.Value
			}
			parts = append(parts, googleGenAI.FunctionResponse{
				Name:     p.FunctionName,
				Response: args,
			})
		}
	}

	resp, err := a.session.SendMessage(ctx, parts...)
	if err != nil {
		return nil, snapmatchai.NewError(err, "failed to send message", 400)
	}
	var aiResp []snapmatchai.AIResponse
	for _, part := range resp.Candidates[0].Content.Parts {
		funcPart, ok := part.(googleGenAI.FunctionCall)
		if ok {
			args := make([]snapmatchai.FunctionArgs, len(funcPart.Args))
			i := 0
			for key, arg := range funcPart.Args {
				args[i] = snapmatchai.FunctionArgs{
					Name:  key,
					Value: arg,
				}
				i++
			}
			log.Println(args)
			aiResp = append(aiResp, snapmatchai.FunctionCall{
				FunctionName: funcPart.Name,
				Args:         args,
			})
			continue
		}
		if text, ok := part.(googleGenAI.Text); ok {
			aiResp = append(aiResp, snapmatchai.Text(text))
			continue
		}
	}
	return aiResp, nil
}

func (c *Client) aiToolsToGenAI(funcs []snapmatchai.AITools) *googleGenAI.Tool {
	var funcDeclarations []*googleGenAI.FunctionDeclaration
	for _, f := range funcs {
		props := make(map[string]*googleGenAI.Schema, len(f.Props))
		for _, p := range f.Props {
			props[p.Key] = &googleGenAI.Schema{Description: p.Description, Type: convertSnapmatchAIType(p.Type)}
		}
		funcDeclr := &googleGenAI.FunctionDeclaration{
			Name:        f.Name,
			Description: f.Description,
			Parameters: &googleGenAI.Schema{
				Type:       googleGenAI.TypeObject,
				Properties: props,
			},
		}
		funcDeclarations = append(funcDeclarations, funcDeclr)
	}

	return &googleGenAI.Tool{
		FunctionDeclarations: funcDeclarations,
	}
}

func convertSnapmatchAIType(s snapmatchai.AIToolPropsType) googleGenAI.Type {
	switch s {
	case snapmatchai.AIToolPropsTypeString:
		return googleGenAI.TypeString
	case snapmatchai.AIToolPropsTypeInt:
		return googleGenAI.TypeInteger
	}
	return googleGenAI.TypeString
}
