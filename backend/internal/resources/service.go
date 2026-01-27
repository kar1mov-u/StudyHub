package resources

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type ResourceRepository interface {
	GetResourceByID(context.Context, uuid.UUID) (Resource, error)
	SaveResource(context.Context, Resource) error
	ListResourcesByWeek(context.Context, uuid.UUID) ([]Resource, error)
}

type ResourceService struct {
	resourceRepo ResourceRepository
	filesStorage FileStorage
}

func NewResourceService(repo ResourceRepository, storage FileStorage) *ResourceService {
	return &ResourceService{resourceRepo: repo, filesStorage: storage}
}

func (s *ResourceService) UploadResource(ctx context.Context, body io.Reader, filename string) error {

	return nil
}
