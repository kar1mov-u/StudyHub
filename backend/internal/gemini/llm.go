package gemini

import (
	"context"
	"io"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
}

func NewGeminiClient(key string) *GeminiClient {
	ctx := context.Background()
	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  key,
		Backend: genai.BackendGeminiAPI,
	})
	return &GeminiClient{client: client}
}

func (gc *GeminiClient) Upload(ctx context.Context, file io.ReadCloser) error {
	parts := []*genai.Part{
		&genai.Part{
			InlineData: &genai.Blob{
				
			}
		}
	}
}
