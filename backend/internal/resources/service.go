package resources

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"log/slog"

	"github.com/google/uuid"
)

type ResourceRepository interface {
	Create(context.Context, Resource) error
	ResourceExists(ctx context.Context, hash string) (uuid.UUID, bool, error)
	// GetResourceByID(context.Context, uuid.UUID) (Resource, error)
	ListResourcesByWeek(context.Context, uuid.UUID) ([]Resource, error)
	CreateUserResource(ctx context.Context, resource Resource) error
	CreateWeekResource(ctx context.Context, resource Resource) error
}

type ResourceService struct {
	resourceRepo ResourceRepository
	filesStorage FileStorage
}

func NewResourceService(repo ResourceRepository, storage FileStorage) *ResourceService {
	return &ResourceService{resourceRepo: repo, filesStorage: storage}
}

func (s *ResourceService) UploadResource(ctx context.Context, body io.Reader, size int64, resource Resource) error {
	log.Println("starting upload srv method")
	hasher := sha256.New()

	// //use UUID for object key in AWS
	storageObjectID := uuid.New()

	tr := io.TeeReader(body, hasher)

	storageObjectUrl, err := s.filesStorage.UploadObject(ctx, storageObjectID.String(), size, tr)
	if err != nil {
		return err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	resourceID, exists, err := s.resourceRepo.ResourceExists(ctx, hash)
	if err != nil {
		return err
	}

	//delete the file from storage if already exists
	if exists {
		slog.Info("resource exists")
		//assign the ID returned from the DB, instead of ID in the request
		resource.ID = resourceID

		// in the backgound delete the object from the storage
		go func() {
			err = s.filesStorage.DeleteObject(context.TODO(), storageObjectID.String())
			if err != nil {
				slog.Info("failed to delete from storage", "err", err)
			} else {
				slog.Info("successfully deleted the file from the storage")
			}
		}()
	}

	//if not exists should create the record on DB
	if !exists {
		resource.Hash = hash
		resource.Url = storageObjectUrl
		err = s.resourceRepo.Create(ctx, resource)
		if err != nil {
			return err
		}
	}

	//
	err = s.resourceRepo.CreateUserResource(ctx, resource)
	if err != nil {
		return err
	}
	return s.resourceRepo.CreateWeekResource(ctx, resource)

}

func (s *ResourceService) ListResourcesForWeek(ctx context.Context, weekID uuid.UUID) ([]Resource, error) {
	return s.resourceRepo.ListResourcesByWeek(ctx, weekID)
}
