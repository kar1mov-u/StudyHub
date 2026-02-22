package resources

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"

	"github.com/google/uuid"
)

var ErrResourceExists = errors.New("resource already exists")

type ResourceRepository interface {
	CreateFileResource(ctx context.Context, resource Resource) error
	CreateLinkResource(ctx context.Context, resource Resource) error
	CreateStorageObject(ctx context.Context, object storageObject) error
	CreateUserResource(ctx context.Context, resource Resource) error
	CreateWeekResource(ctx context.Context, resource Resource) error
	ObjectExists(ctx context.Context, hash string) (uuid.UUID, bool, error)
	ListResourcesByWeek(ctx context.Context, weekID uuid.UUID) ([]ResourceWithUser, error)
	ListUserResources(ctx context.Context, userID uuid.UUID) ([]UserResources, error)
	LinkExistsInWeek(ctx context.Context, resource Resource) (bool, error)
	FileExistsInWeek(ctx context.Context, hash string, weekID uuid.UUID) (bool, error)
	ListOrphanObjects(ctx context.Context) ([]uuid.UUID, error)
	DeleteStorageObjecst(ctx context.Context, ids []uuid.UUID) error
	DeleteResource(ctx context.Context, userID, resourceID uuid.UUID) error
}

type Queue interface {
	Publish(ctx context.Context, objectID uuid.UUID) error
}

type ResourceService struct {
	resourceRepo ResourceRepository
	filesStorage FileStorage
	queue        Queue
}

func NewResourceService(repo ResourceRepository, storage FileStorage, queue Queue) *ResourceService {
	return &ResourceService{resourceRepo: repo, filesStorage: storage, queue: queue}
}

func (s *ResourceService) UploadResource(ctx context.Context, body io.Reader, size int64, resource Resource) error {
	hasher := sha256.New()

	//use UUID for object key in AWS
	storageObjectID := uuid.New()
	//teeReader to read the file stream to both hahser and the cloud storage
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
		resource.ObjectID = &storageObjectID

		storageObject := storageObject{ID: storageObjectID}
		storageObject.Hash = hash
		storageObject.URL = storageObjectUrl
		err = s.resourceRepo.CreateStorageObject(ctx, storageObject)
		if err != nil {
			return err
		}
		//here should upload to the queue
		err = s.queue.Publish(ctx, storageObject.ID)
		if err != nil {
			slog.Error("failed to publish message", "err", err.Error())
		}
		slog.Info("published message")
	}

	//check if it exists in the week, to prevent resource deduplication
	exists, err = s.resourceRepo.FileExistsInWeek(ctx, hash, resource.WeekID)
	if err != nil {
		return err
	}
	if exists {
		return ErrResourceExists
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
	//first check if link exists in that week
	exists, err := s.resourceRepo.LinkExistsInWeek(ctx, resource)
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

func (s *ResourceService) ListResourcesForWeek(ctx context.Context, weekID uuid.UUID) ([]ResourceWithUser, error) {
	return s.resourceRepo.ListResourcesByWeek(ctx, weekID)
}

func (s *ResourceService) ListResourceForUser(ctx context.Context, userID uuid.UUID) ([]UserResources, error) {
	return s.resourceRepo.ListUserResources(ctx, userID)
}

func (s *ResourceService) GetResource(ctx context.Context, id uuid.UUID) (string, error) {
	//first get the key for the object
	return s.filesStorage.CreatePresingedURL(ctx, id.String())

}

func (s *ResourceService) CleanOrphanObjects(ctx context.Context) ([]uuid.UUID, error) {
	//get the keys of the files that are not referenced.
	ids, err := s.resourceRepo.ListOrphanObjects(ctx)
	if err != nil {
		return []uuid.UUID{}, err
	}

	deletedIds := make([]uuid.UUID, 0)
	for _, id := range ids {
		err = s.filesStorage.DeleteObject(context.TODO(), id.String())
		if err != nil {
			slog.Error("failed to delete storage Object", "err", err)
		} else {
			deletedIds = append(deletedIds, id)
		}
	}
	err = s.resourceRepo.DeleteStorageObjecst(ctx, deletedIds)
	if err != nil {
		slog.Error("failed to delete storage_objects from DB, but deleted from the storage", "err", err, "ids", ids)
		return []uuid.UUID{}, err
	}
	return deletedIds, nil
}

func (s *ResourceService) DeleteResource(ctx context.Context, userID, resourceID uuid.UUID) error {
	return s.resourceRepo.DeleteResource(ctx, userID, resourceID)
}
