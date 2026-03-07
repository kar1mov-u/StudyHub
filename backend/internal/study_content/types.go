package studycontent

import "github.com/google/uuid"

type Flashcard struct {
	ObjectID *uuid.UUID
	UserID   *uuid.UUID
	WeekID   *uuid.UUID
	Front    string
	Back     string
}
