package content

import "github.com/google/uuid"

type Flashcard struct {
	ID       uuid.UUID
	ObjectID *uuid.UUID
	UserID   *uuid.UUID
	WeekID   *uuid.UUID
	Front    string
	Back     string
}
