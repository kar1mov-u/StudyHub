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
	ObjectID     *uuid.UUID
	ExternalLink *string
	ResourceType ResourceType
	Name         string
	CreatedAt    time.Time
}

type storageObject struct {
	ID   uuid.UUID
	Hash string
	URL  string
}
