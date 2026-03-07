package gemini

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

const prompt = `Prompt:
You are an expert educator and data extraction assistant. Your task is to read the provided document text and generate a comprehensive set of high-quality flashcards.

Instructions:
Content: Identify key concepts, definitions, dates, and relationships. Create "Front" (Question/Term) and "Back" (Answer/Definition) pairs.

Atomicity: Each flashcard should cover exactly one discrete idea to ensure effective active recall.

Format: Your entire response must be a single, valid JSON object containing an array of objects. Do not include any introductory or concluding text.
Give just 5 flashcards
Required JSON Schema:
JSON
{
 [
    {
      "front": "The question or term goes here",
      "back": "The concise answer or definition goes here"
    }
  ]
}`

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

func (gc *GeminiClient) GenerateFlashCards(ctx context.Context, file io.ReadCloser) (string, error) {

	defer file.Close()
	uploadConfig := &genai.UploadFileConfig{MIMEType: "application/pdf"}
	uploadedFile, err := gc.client.Files.Upload(ctx, file, uploadConfig)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to the Gemini: %w", err)
	}

	promptParts := []*genai.Part{
		genai.NewPartFromURI(uploadedFile.URI, uploadConfig.MIMEType),
		genai.NewPartFromText(prompt),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(promptParts, genai.RoleUser),
	}

	result, err := gc.client.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview",
		contents,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}
