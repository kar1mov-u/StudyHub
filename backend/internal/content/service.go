package content

import (
	"context"
	"errors"
	"io"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	isPdf(ctx context.Context, id uuid.UUID) string
	CreateCardsFromObject(ctx context.Context, cards []Flashcard) error
	ListCardsFromObjects(ctx context.Context, ids []uuid.UUID) ([]Flashcard, error)

	// User Deck Methods
	AddCardToUserDeck(ctx context.Context, userID, weekID, flashcardID uuid.UUID) error
	CreateCustomCardInDeck(ctx context.Context, userID, weekID uuid.UUID, front, back string) (UserDeckCard, error)
	RemoveCardFromUserDeck(ctx context.Context, cardID, userID uuid.UUID) error
	GetUserDeckForWeek(ctx context.Context, userID, weekID uuid.UUID) ([]UserDeckCard, error)
	GetUserDeckCard(ctx context.Context, cardID, userID uuid.UUID) (UserDeckCard, error)
	UpdateUserDeckCard(ctx context.Context, cardID, userID uuid.UUID, front, back *string) error
	RecordCardReview(ctx context.Context, cardID, userID uuid.UUID, difficultyRating int) error
	GetDeckStatistics(ctx context.Context, userID, weekID uuid.UUID) (DeckStats, error)
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

// AddCardToUserDeck adds an auto-generated flashcard to user's personal deck
func (s *ContentService) AddCardToUserDeck(ctx context.Context, userID, weekID, flashcardID uuid.UUID) error {
	// Validate that the flashcard exists
	cards, err := s.contentRepository.ListCardsFromObjects(ctx, []uuid.UUID{flashcardID})
	if err != nil {
		return err
	}
	if len(cards) == 0 {
		return errors.New("flashcard not found")
	}

	// Add to user's deck (duplicate check handled by DB unique constraint)
	return s.contentRepository.AddCardToUserDeck(ctx, userID, weekID, flashcardID)
}

// CreateCustomCard creates a custom flashcard directly in user's deck
func (s *ContentService) CreateCustomCard(ctx context.Context, userID, weekID uuid.UUID, front, back string) (UserDeckCard, error) {
	if front == "" || back == "" {
		return UserDeckCard{}, errors.New("front and back cannot be empty")
	}
	return s.contentRepository.CreateCustomCardInDeck(ctx, userID, weekID, front, back)
}

// GetUserDeckForWeek retrieves all cards in user's deck for a specific week
func (s *ContentService) GetUserDeckForWeek(ctx context.Context, userID, weekID uuid.UUID) ([]UserDeckCard, error) {
	return s.contentRepository.GetUserDeckForWeek(ctx, userID, weekID)
}

// UpdateDeckCard updates card content (only if user owns it)
func (s *ContentService) UpdateDeckCard(ctx context.Context, cardID, userID uuid.UUID, front, back *string) error {
	// Verify the card exists and belongs to the user
	_, err := s.contentRepository.GetUserDeckCard(ctx, cardID, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("card not found or access denied")
		}
		return err
	}

	return s.contentRepository.UpdateUserDeckCard(ctx, cardID, userID, front, back)
}

// RemoveCardFromDeck removes a card from user's deck
func (s *ContentService) RemoveCardFromDeck(ctx context.Context, cardID, userID uuid.UUID) error {
	err := s.contentRepository.RemoveCardFromUserDeck(ctx, cardID, userID)
	if err == pgx.ErrNoRows {
		return errors.New("card not found or access denied")
	}
	return err
}

// RecordCardReview records a study session for a card
func (s *ContentService) RecordCardReview(ctx context.Context, cardID, userID uuid.UUID, difficultyRating int) error {
	// Validate difficulty rating
	if difficultyRating < 1 || difficultyRating > 5 {
		return errors.New("difficulty rating must be between 1 and 5")
	}

	err := s.contentRepository.RecordCardReview(ctx, cardID, userID, difficultyRating)
	if err == pgx.ErrNoRows {
		return errors.New("card not found or access denied")
	}
	return err
}

// GetDeckStats retrieves statistics for a user's deck in a specific week
func (s *ContentService) GetDeckStats(ctx context.Context, userID, weekID uuid.UUID) (DeckStats, error) {
	return s.contentRepository.GetDeckStatistics(ctx, userID, weekID)
}
