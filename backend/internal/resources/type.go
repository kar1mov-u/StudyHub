package resources

import (
	"time"

	"github.com/google/uuid"
)

type ResourceType string

const (
	ResourceFile ResourceType = "file"
	ResourceLink ResourceType = "link"
	ResourceNote ResourceType = "note"
)

type Resource struct {
	ID           uuid.UUID
	WeekID       uuid.UUID
	UserID       uuid.UUID
	ResourceType ResourceType
	Hash         string
	Name         string
	Url          string //url in the cloud storage
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
