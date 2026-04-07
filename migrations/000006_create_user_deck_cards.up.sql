CREATE TABLE IF NOT EXISTS user_deck_cards (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_id UUID NOT NULL REFERENCES weeks(id) ON DELETE CASCADE,
    source_flashcard_id UUID REFERENCES flashcards(id) ON DELETE SET NULL,
    front TEXT NOT NULL,
    back TEXT NOT NULL,
    is_custom BOOLEAN NOT NULL DEFAULT FALSE,
    last_reviewed_at TIMESTAMP,
    review_count INT NOT NULL DEFAULT 0,
    difficulty_rating INT CHECK (difficulty_rating BETWEEN 1 AND 5),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, week_id, source_flashcard_id)
);

CREATE INDEX idx_user_deck_cards_user_week ON user_deck_cards(user_id, week_id);
CREATE INDEX idx_user_deck_cards_source ON user_deck_cards(source_flashcard_id);
