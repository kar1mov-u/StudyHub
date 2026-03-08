package content

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

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
