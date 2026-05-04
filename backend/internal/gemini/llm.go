package gemini

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

const chatSystemContext = `You are a helpful assistant for StudyHub, an academic study management platform.

StudyHub helps students manage their academic modules and study resources. Here is what the app does:

- Modules: Academic courses organised by department (e.g. "Introduction to Computer Science"). Each module has one or more runs.
- Module Runs: A specific instance of a module in a semester and year (e.g. Spring 2024). Each run contains weekly sessions.
- Weeks: Individual weeks within a module run, each holding uploaded resources.
- Resources: Files (PDFs, docs, etc.) or external links that students upload to a specific week. Files are stored in AWS S3.
- Flashcard Decks: AI-generated flashcards from uploaded PDF files that help students study. Users can also create custom cards.
- Comments: Students can comment on weekly resources and upvote/downvote others' comments.
- User Profiles: Each user can view their own uploads and those of other students.
- Academic Terms: Admin-managed terms (semester + year) that group module runs.

Navigation structure: Modules → Module Detail (runs & weeks) → Week Detail (resources, comments, flashcards)

Answer questions about how to use the app, its features, and how to navigate it. If asked about real-time data (e.g. "what modules exist"), explain you can't access live data but describe where to find it in the app.`

const prompt = `Prompt:
You are an expert educator and data extraction assistant. Your task is to read the provided document text and generate a comprehensive set of high-quality flashcards.
Instructions:
Content: Identify key concepts, definitions, dates, and relationships. Create "Front" (Question/Term) and "Back" (Answer/Definition) pairs.
Atomicity: Each flashcard should cover exactly one discrete idea to ensure effective active recall.
Format: Your entire response must be a single, valid JSON object containing an array of objects. Do not include any introductory or concluding text.
Required JSON Schema:
{
 [
    {
      "front": "The question or term goes here",
      "back": "The concise answer or definition goes here"
    }
  ]
}
Output shouldn't include with starting and trailing JSON markdown, and do not include \n whitespaces  
`

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

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("failed to close file: %v\n", closeErr)
		}
	}()
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
		"gemini-2.5-flash",
		contents,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}

func (gc *GeminiClient) Chat(ctx context.Context, message string) (string, error) {
	fullPrompt := chatSystemContext + "\n\nUser question: " + message

	contents := []*genai.Content{
		genai.NewContentFromParts([]*genai.Part{genai.NewPartFromText(fullPrompt)}, genai.RoleUser),
	}

	result, err := gc.client.Models.GenerateContent(ctx, "gemini-2.5-flash-lite", contents, nil)
	if err != nil {
		return "", fmt.Errorf("gemini chat failed: %w", err)
	}
	return result.Text(), nil
}
