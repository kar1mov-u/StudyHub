package gemini

import (
	"context"
	"fmt"
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

func (gc *GeminiClient) Generate(ctx context.Context, file io.ReadCloser) (string, error) {
	uploadConfig := &genai.UploadFileConfig{MIMEType: "application/pdf"}
	uploadedFile, err := gc.client.Files.Upload(ctx, file, uploadConfig)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to the Gemini: %w", err)
	}

	promptParts := []*genai.Part{
		genai.NewPartFromURI(uploadedFile.URI, uploadConfig.MIMEType),
		genai.NewPartFromText("Summarize this document in 2 words"),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(promptParts, genai.RoleUser),
	}

	result, _ := gc.client.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview",
		contents,
		nil,
	)

	return result.Text(), nil
}
