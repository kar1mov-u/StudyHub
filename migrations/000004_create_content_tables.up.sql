CREATE TABLE IF NOT EXISTS flashcards(
    id UUID PRIMARY KEY,
    storage_object_id UUID REFERENCES storage_objects(id) ON DELETE CASCADE, 
    user_id UUID REFERENCES users(id),
    week_id UUID REFERENCES weeks(id),
    front TEXT NOT NULL, 
    back TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
)

