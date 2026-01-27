package http

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plaintext password with a bcrypt hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func parseUUID(w http.ResponseWriter, param string) (uuid.UUID, bool) {
	id, err := uuid.Parse(param)
	if err != nil {
		ResponseWithErr(w, http.StatusBadRequest, "invalid UUID format")
		return uuid.Nil, false
	}
	return id, true
}

func getUserID(r *http.Request) string {
	id := r.Context().Value("userID").(string)
	return id
}

func isNotFoundError(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
