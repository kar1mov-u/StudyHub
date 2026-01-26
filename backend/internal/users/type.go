package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
	IsAdmin   bool
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
