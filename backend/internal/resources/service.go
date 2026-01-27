package resources

import (
	"context"

	"github.com/google/uuid"
)

type ResourceRepository interface {
	Create(context.Context, Resource) error
	ResourceExists(ctx context.Context, hash string) (uuid.UUID, bool, error)
	GetResourceByID(context.Context, uuid.UUID) (Resource, error)
	ListResourcesByWeek(context.Context, uuid.UUID) ([]Resource, error)
}

type ResourceService struct {
	resourceRepo ResourceRepository
	filesStorage FileStorage
}

func NewResourceService(repo ResourceRepository, storage FileStorage) *ResourceService {
	return &ResourceService{resourceRepo: repo, filesStorage: storage}
}

// func (s *ResourceService) UploadResource(ctx context.Context, body io.Reader, resource Resource) error {
// 	hasher := sha256.New()
// 	//need to create pipe, becuse S3 storage needs the io.reader
// 	pr, pw := io.Pipe()
// 	go func() {
// 		mw := io.MultiWriter(hasher, pw)
// 		defer pr.Close()

// 		_, err := io.Copy(mw, body)
// 		if err != nil {
// 			slog.Error("failed to copy file to multiWriter", "err", err.Error())
// 		}
// 	}()

// 	err := s.filesStorage.UploadFile(ctx, resource.Name, pr)
// 	if err != nil {
// 		return err
// 	}

// 	hash := hex.EncodeToString(hasher.Sum(nil))

// 	resourceID, exists, err := s.resourceRepo.ResourceExists(ctx, hash)
// 	if err != nil {
// 		return err
// 	}

// }
