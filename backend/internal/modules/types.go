package internal

import (
	"time"

	"github.com/google/uuid"
)

type ModulePage struct {
	Module Module
	Run    ModuleRun
	Weeks  []Week
}

type Module struct {
	ID             uuid.UUID
	Code           string
	Name           string
	Description    string
	DepartmentName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ModuleRun struct {
	ID        uuid.UUID
	ModuleID  uuid.UUID
	Year      int
	Semester  string
	Weeks     int
	IsActive  bool
	CreatedAt time.Time
}

type Week struct {
	ID          uuid.UUID
	ModuleRunID uuid.UUID
	Number      int
	Topic       string
}
