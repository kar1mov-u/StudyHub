package content

import (
	"time"

	"github.com/google/uuid"
)

type Flashcard struct {
	ID       uuid.UUID
	ObjectID *uuid.UUID
	UserID   *uuid.UUID
	WeekID   *uuid.UUID
	Front    string
	Back     string
}

// UserDeckCard represents a flashcard in a user's personal deck for a specific week
type UserDeckCard struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	WeekID            uuid.UUID
	SourceFlashcardID *uuid.UUID // nil for custom cards
	Front             string
	Back              string
	IsCustom          bool
	LastReviewedAt    *time.Time
	ReviewCount       int
	DifficultyRating  *int // 1-5 scale
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// AddCardToDeckRequest for adding auto-generated card to user deck
type AddCardToDeckRequest struct {
	FlashcardID string `json:"flashcard_id"`
}

// CreateCustomCardRequest for creating custom flashcard
type CreateCustomCardRequest struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

// UpdateCardRequest for updating card content
type UpdateCardRequest struct {
	Front *string `json:"front,omitempty"`
	Back  *string `json:"back,omitempty"`
}

// RecordReviewRequest for recording card review
type RecordReviewRequest struct {
	DifficultyRating int `json:"difficulty_rating"` // 1-5
}

// DeckStats represents statistics for a user's deck in a week
type DeckStats struct {
	TotalCards     int        `json:"total_cards"`
	ReviewedCards  int        `json:"reviewed_cards"`
	AverageRating  float64    `json:"average_rating"`
	LastReviewedAt *time.Time `json:"last_reviewed_at,omitempty"`
}
