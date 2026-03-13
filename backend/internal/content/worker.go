package content

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
)

const GOTENBERG_URL = "http://gotenberg:3000/forms/libreoffice/convert"

func (s *ContentService) startWorkers() {
	for i := range 5 {
		slog.Info("started worker", "id", i)
		go s.worker()
	}
}

func (s *ContentService) worker() {
	for msg := range s.delivery {

		key := string(msg.Body)
		slog.Info("started on job with object id", "ID", string(msg.Body))

		file, err := s.fileStorage.GetObject(context.Background(), key)
		if err != nil {
			slog.Error("error getting file from storage", "err", err)
			continue
		}

		//if file is not pdf, we get the file type and send to the ai service to convert it to pdf, then we continue with the pdf file
		file, err = s.convertToPdf(context.TODO(), key, file)
		if err != nil {
			slog.Error("falied to convert pdf", "err", err)
			continue
		}

		result, err := s.ai.GenerateFlashCards(context.Background(), file)
		if err != nil {
			slog.Error("failed to generate content", "error: ", err)
			continue
		}
		flashcards, err := cleanupResult(result, string(msg.Body))
		if err != nil {
			slog.Error("failed to clenup flashcards", "error: ", err)
			continue
		}

		//save to the DB

		err = s.contentRepository.CreateCardsFromObject(context.Background(), flashcards)

		slog.Info("finished job ", " object id :", key)

	}
	slog.Info("worker exiting")
}

func cleanupResult(data, id string) ([]Flashcard, error) {
	var res []Flashcard
	err := json.Unmarshal([]byte(data), &res)
	if err != nil {
		return []Flashcard{}, fmt.Errorf("%w | Payload: %s", err, data)
	}
	objectID, _ := uuid.Parse(id)
	for i := range res {
		res[i].ID = uuid.New()
		res[i].ObjectID = &objectID
	}

	return res, nil
}

func (s *ContentService) convertToPdf(ctx context.Context, idString string, body io.ReadCloser) (io.ReadCloser, error) {
	id, _ := uuid.Parse(idString)
	if s.contentRepository.isPdf(ctx, id) {
		return body, nil
	}

	return makeGotenbergRequest(ctx, body, idString)

}

func makeGotenbergRequest(ctx context.Context, body io.Reader, name string) (io.ReadCloser, error) {
	//everything read from pr will be writen to pw
	slog.Info("starting pdf conversion")

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer writer.Close()

		part, err := writer.CreateFormFile("file", name)
		if err != nil {
			return
		}

		_, _ = io.Copy(part, body)
	}()
	req, err := http.NewRequest("POST", GOTENBERG_URL, pr)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}
	return resp.Body, err

}
