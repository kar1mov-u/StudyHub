package content

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContentRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewContentRepositoryPostgres(p *pgxpool.Pool) *ContentRepositoryPostgres {
	return &ContentRepositoryPostgres{
		pool: p,
	}
}

// should do the batch insert
func (r *ContentRepositoryPostgres) CreateCardsFromObject(ctx context.Context, cards []Flashcard) error {
	query := `INSERT INTO flashcards(id, storage_object_id,front, back) VALUES ($1, $2, $3, $4)`

	batch := pgx.Batch{}
	for _, card := range cards {
		batch.Queue(query, card.ID, card.ObjectID, card.Front, card.Back)
	}
	err := r.pool.SendBatch(ctx, &batch).Close()
	return err
}
func (r *ContentRepositoryPostgres) ListCardsFromObjects(ctx context.Context, ids []uuid.UUID) ([]Flashcard, error) {

	query := `SELECT id, storage_object_id, front, back FROM flashcards WHERE storage_object_id = ANY ($1)`
	cards := make([]Flashcard, 0)
	rows, err := r.pool.Query(ctx, query, ids)
	if err != nil {
		return []Flashcard{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var card Flashcard
		err = rows.Scan(&card.ID, &card.ObjectID, &card.Front, &card.Back)
		if err != nil {
			return []Flashcard{}, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (r *ContentRepositoryPostgres) isPdf(ctx context.Context, id uuid.UUID) string {
	query := `SELECT file_type from storage_objects WHERE id=$1`
	var fileType string
	err := r.pool.QueryRow(ctx, query, id).Scan(&fileType)
	if err != nil {
		return "pdf"
	}
	return fileType

}

// AddCardToUserDeck adds an auto-generated flashcard to user's personal deck
func (r *ContentRepositoryPostgres) AddCardToUserDeck(ctx context.Context, userID, weekID, flashcardID uuid.UUID) error {
	// First, get the card content from the flashcards table
	var front, back string
	getCardQuery := `SELECT front, back FROM flashcards WHERE id = $1`
	err := r.pool.QueryRow(ctx, getCardQuery, flashcardID).Scan(&front, &back)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("flashcard not found")
		}
		return err
	}

	// Insert into user_deck_cards
	insertQuery := `
		INSERT INTO user_deck_cards (id, user_id, week_id, source_flashcard_id, front, back, is_custom)
		VALUES ($1, $2, $3, $4, $5, $6, false)
	`
	_, err = r.pool.Exec(ctx, insertQuery, uuid.New(), userID, weekID, flashcardID, front, back)
	return err
}

// CreateCustomCardInDeck creates a custom flashcard in user's deck
func (r *ContentRepositoryPostgres) CreateCustomCardInDeck(ctx context.Context, userID, weekID uuid.UUID, front, back string) (UserDeckCard, error) {
	card := UserDeckCard{
		ID:        uuid.New(),
		UserID:    userID,
		WeekID:    weekID,
		Front:     front,
		Back:      back,
		IsCustom:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO user_deck_cards (id, user_id, week_id, source_flashcard_id, front, back, is_custom)
		VALUES ($1, $2, $3, NULL, $4, $5, true)
		RETURNING created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query, card.ID, card.UserID, card.WeekID, card.Front, card.Back).
		Scan(&card.CreatedAt, &card.UpdatedAt)
	if err != nil {
		return UserDeckCard{}, err
	}

	return card, nil
}

// RemoveCardFromUserDeck removes a card from user's deck
func (r *ContentRepositoryPostgres) RemoveCardFromUserDeck(ctx context.Context, cardID, userID uuid.UUID) error {
	query := `DELETE FROM user_deck_cards WHERE id = $1 AND user_id = $2`
	result, err := r.pool.Exec(ctx, query, cardID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// GetUserDeckForWeek retrieves all cards in user's deck for a specific week
func (r *ContentRepositoryPostgres) GetUserDeckForWeek(ctx context.Context, userID, weekID uuid.UUID) ([]UserDeckCard, error) {
	query := `
		SELECT id, user_id, week_id, source_flashcard_id, front, back, is_custom, 
		       last_reviewed_at, review_count, difficulty_rating, created_at, updated_at
		FROM user_deck_cards
		WHERE user_id = $1 AND week_id = $2
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID, weekID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make([]UserDeckCard, 0)
	for rows.Next() {
		var card UserDeckCard
		err := rows.Scan(
			&card.ID, &card.UserID, &card.WeekID, &card.SourceFlashcardID,
			&card.Front, &card.Back, &card.IsCustom, &card.LastReviewedAt,
			&card.ReviewCount, &card.DifficultyRating, &card.CreatedAt, &card.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, rows.Err()
}

// GetUserDeckCard retrieves a single card from user's deck
func (r *ContentRepositoryPostgres) GetUserDeckCard(ctx context.Context, cardID, userID uuid.UUID) (UserDeckCard, error) {
	query := `
		SELECT id, user_id, week_id, source_flashcard_id, front, back, is_custom, 
		       last_reviewed_at, review_count, difficulty_rating, created_at, updated_at
		FROM user_deck_cards
		WHERE id = $1 AND user_id = $2
	`
	var card UserDeckCard
	err := r.pool.QueryRow(ctx, query, cardID, userID).Scan(
		&card.ID, &card.UserID, &card.WeekID, &card.SourceFlashcardID,
		&card.Front, &card.Back, &card.IsCustom, &card.LastReviewedAt,
		&card.ReviewCount, &card.DifficultyRating, &card.CreatedAt, &card.UpdatedAt,
	)
	if err != nil {
		return UserDeckCard{}, err
	}
	return card, nil
}

// UpdateUserDeckCard updates the content of a card in user's deck
func (r *ContentRepositoryPostgres) UpdateUserDeckCard(ctx context.Context, cardID, userID uuid.UUID, front, back *string) error {
	// Build dynamic query based on which fields are being updated
	if front == nil && back == nil {
		return nil // Nothing to update
	}

	query := `UPDATE user_deck_cards SET updated_at = NOW()`
	args := []interface{}{}
	argPos := 1

	if front != nil {
		query += fmt.Sprintf(", front = $%d", argPos)
		args = append(args, *front)
		argPos++
	}
	if back != nil {
		query += fmt.Sprintf(", back = $%d", argPos)
		args = append(args, *back)
		argPos++
	}

	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", argPos, argPos+1)
	args = append(args, cardID, userID)

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// RecordCardReview records a review session for a card
func (r *ContentRepositoryPostgres) RecordCardReview(ctx context.Context, cardID, userID uuid.UUID, difficultyRating int) error {
	query := `
		UPDATE user_deck_cards 
		SET last_reviewed_at = NOW(),
		    review_count = review_count + 1,
		    difficulty_rating = $1,
		    updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`
	result, err := r.pool.Exec(ctx, query, difficultyRating, cardID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// GetDeckStatistics retrieves statistics for a user's deck in a specific week
func (r *ContentRepositoryPostgres) GetDeckStatistics(ctx context.Context, userID, weekID uuid.UUID) (DeckStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_cards,
			COUNT(last_reviewed_at) as reviewed_cards,
			COALESCE(AVG(difficulty_rating), 0) as avg_rating,
			MAX(last_reviewed_at) as last_reviewed
		FROM user_deck_cards
		WHERE user_id = $1 AND week_id = $2
	`
	var stats DeckStats
	err := r.pool.QueryRow(ctx, query, userID, weekID).Scan(
		&stats.TotalCards,
		&stats.ReviewedCards,
		&stats.AverageRating,
		&stats.LastReviewedAt,
	)
	if err != nil {
		return DeckStats{}, err
	}
	return stats, nil
}
