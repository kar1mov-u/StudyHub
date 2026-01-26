package auth

import (
	"StudyHub/backend/internal/users"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtKey   string
	userRepo users.UserRepository
}

func NewAuthSerivce(key string, userRep users.UserRepository) *AuthService {
	return &AuthService{
		jwtKey:   key,
		userRepo: userRep,
	}
}

// this func will handle the login functionality, will return the JWT
func (s *AuthService) LoginUser(ctx context.Context, email, password string) (string, error) {

	//get the user from the repo
	userID, hash, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	//check for the password
	if !CheckPasswordHash(password, hash) {
		return "", fmt.Errorf("invalid password was provided")
	}

	//create a JWT

	jwt, err := createNewJWT(userID.String(), s.jwtKey)
	if err != nil {
		slog.Error("shit")
		return "", err
	}

	return jwt, nil
}

func (s *AuthService) IsAdmin(ctx context.Context, id uuid.UUID) bool {
	ok, err := s.userRepo.IsAdmin(ctx, id)
	if err != nil || !ok {
		return false
	}
	return true
}

// CheckPasswordHash compares a plaintext password with a bcrypt hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
