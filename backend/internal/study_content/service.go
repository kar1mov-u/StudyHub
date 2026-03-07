package studycontent

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Consume() chan amqp.Delivery
}

type FileStorage interface {
	GetObject(ctx context.Context, objectID string) (io.ReadCloser, error)
}

type AI interface {
	GenerateFlashCards(ctx context.Context, file io.ReadCloser) (string, error)
}

type StudyContentRepository interface {
}

type StudyContentService struct {
	contentRepository StudyContentRepository
	queue             Queue
	delivery          chan amqp.Delivery
	fileStorage       FileStorage
	ai                AI
}

func NewStudyContentService(q Queue, fileStorage FileStorage, ai AI) *StudyContentService {
	s := &StudyContentService{
		queue:       q,
		fileStorage: fileStorage,
		ai:          ai,
		delivery:    q.Consume(),
	}
	s.startWorkers()
	return s
}

func (s *StudyContentService) startWorkers() {
	for i := range 5 {
		slog.Info("started worker", "id", i)
		go s.worker()
	}
}

func (s *StudyContentService) worker() {
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
		//design the DB schema

		log.Println(flashcards)

		slog.Info("finished job ", " object id :", key, "result: ", result)

	}
	slog.Info("worker exiting")
}

func cleanupResult(data, id string) ([]Flashcard, error) {
	var res []Flashcard
	err := json.Unmarshal([]byte(data), &res)
	if err != nil {
		return []Flashcard{}, err
	}
	uuid, _ := uuid.Parse(id)
	for i := range res {
		res[i].ObjectID = &uuid
	}

	return res, nil
}
