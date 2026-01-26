package users

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(context.Context, User) error
	GetByID(context.Context, uuid.UUID) (User, error)
	GetByEmail(context.Context, string) (uuid.UUID, string, error)
	Delete(context.Context, uuid.UUID) error
	List(context.Context) ([]User, error)
	IsAdmin(context.Context, uuid.UUID) (bool, error)
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{userRepo: r}
}

func (s *UserService) Create(ctx context.Context, user User) error {
	return s.userRepo.Create(ctx, user)
}

func (s *UserService) Get(ctx context.Context, id uuid.UUID) (User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) List(ctx context.Context) ([]User, error) {
	return s.userRepo.List(ctx)
}
