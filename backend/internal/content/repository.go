package content

import (
	"context"

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
func (r *ContentRepositoryPostgres) test() ([]Flashcard, error) {
	return []Flashcard{}, nil
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
