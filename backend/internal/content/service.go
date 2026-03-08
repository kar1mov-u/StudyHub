package content

import (
	"context"
	"io"

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

type ContentRepository interface {
	CreateCardsFromObject(ctx context.Context, cards []Flashcard) error
	ListCardsFromObjects(ctx context.Context, ids []uuid.UUID) ([]Flashcard, error)
}

type ContentService struct {
	contentRepository ContentRepository
	queue             Queue
	delivery          chan amqp.Delivery
	fileStorage       FileStorage
	ai                AI
}

func NewContentService(contentRepo ContentRepository, q Queue, fileStorage FileStorage, ai AI) *ContentService {
	s := &ContentService{
		contentRepository: contentRepo,
		queue:             q,
		fileStorage:       fileStorage,
		ai:                ai,
		delivery:          q.Consume(),
	}
	s.startWorkers()
	return s
}

func (s *ContentService) ListSelectedObjectsCards(ctx context.Context, ids []uuid.UUID) ([]Flashcard, error) {
	//for each object get the cards
	cards, err := s.contentRepository.ListCardsFromObjects(ctx, ids)
	if err != nil {
		return []Flashcard{}, err
	}

	return cards, nil

}
