package resources

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"log/slog"

	"github.com/google/uuid"
)

type ResourceRepository interface {
	CreateFileResource(ctx context.Context, resource Resource) error
	CreateStorageObject(ctx context.Context, object storageObject) error
	CreateUserResource(ctx context.Context, resource Resource) error
	CreateWeekResource(ctx context.Context, resource Resource) error
	ObjectExists(ctx context.Context, hash string) (uuid.UUID, bool, error)
	ListResourcesByWeek(ctx context.Context, weekID uuid.UUID) ([]Resource, error)
	LinkExists(ctx context.Context, resource Resource) (bool, error)
	CreateLinkResource(ctx context.Context, resource Resource) error
	// GetStorageKeyForResource(ctx context.Context, resourceID uuid.UUID) (string, error)
}

var ErrResourceExists = errors.New("resource already exists")

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

	objectID, exists, err := s.resourceRepo.ObjectExists(ctx, hash)
	if err != nil {
		return err
	}

	//delete the file from storage if already exists
	if exists {
		slog.Info("resource exists")
		//assign the ID returned from the DB, instead of ID in the request
		resource.ObjectID = &objectID

		// in the backgound delete the object from the storage
		go func() {
			err = s.filesStorage.DeleteObject(context.TODO(), storageObjectID.String())
			if err != nil {
				slog.Info("failed to delete from storage", "err", err)
			} else {
				slog.Info("successfully deleted the file from the storage")
			}
		}()
	} else {
		storageObject := storageObject{ID: storageObjectID}
		resource.ObjectID = &storageObject.ID
		storageObject.Hash = hash
		storageObject.URL = storageObjectUrl
		err = s.resourceRepo.CreateStorageObject(ctx, storageObject)
		if err != nil {
			return err
		}
	}

	err = s.resourceRepo.CreateFileResource(ctx, resource)
	if err != nil {
		return err
	}
	err = s.resourceRepo.CreateUserResource(ctx, resource)
	if err != nil {
		return err
	}
	return s.resourceRepo.CreateWeekResource(ctx, resource)

}

func (s *ResourceService) CreateLinkResource(ctx context.Context, resource Resource) error {
	//first check if link exists
	exists, err := s.resourceRepo.LinkExists(ctx, resource)
	if err != nil {
		return err
	}
	if exists {
		return ErrResourceExists
	}
	err = s.resourceRepo.CreateLinkResource(ctx, resource)
	if err != nil {
		return err
	}
	err = s.resourceRepo.CreateUserResource(ctx, resource)
	if err != nil {
		return err
	}
	return s.resourceRepo.CreateWeekResource(ctx, resource)

}

func (s *ResourceService) ListResourcesForWeek(ctx context.Context, weekID uuid.UUID) ([]Resource, error) {
	return s.resourceRepo.ListResourcesByWeek(ctx, weekID)
}

func (s *ResourceService) GetResource(ctx context.Context, id uuid.UUID) (string, error) {
	//first get the key for the object
	return s.filesStorage.CreatePresingedURL(ctx, id.String())

}
