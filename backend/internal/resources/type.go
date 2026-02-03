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

// this is not Domain type , but instaed a struct that is used when we want +info about the Owner
type ResourceWithUser struct {
	ID           uuid.UUID
	WeekID       uuid.UUID
	UserID       uuid.UUID
	UserName     string
	ObjectID     *uuid.UUID
	ExternalLink *string
	ResourceType ResourceType
	Name         string
	CreatedAt    time.Time
}

// its used to get resources that are shared by user, should include info about when and where uploaded(Module,Semester, WeekNum)
type UserResources struct {
	ID           uuid.UUID
	WeekID       uuid.UUID // we should give uuid for that week, so by clicking on link, user will be redirected to the weeks page
	UserID       uuid.UUID
	ModuleName   string //for better context
	Semester     string
	Year         int
	WeekNumber   int
	ObjectID     *uuid.UUID // do we need this, we shouldnt be able to download from there, should be just for listing, later can change
	ExternalLink *string
	ResourceType ResourceType
	Name         string
	CreatedAt    time.Time
}
