package http

import (
	"StudyHub/internal/users"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserDTO struct {
	ID        uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
}

func (s *HTTPServer) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userReqeust CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&userReqeust)
	if err != nil {
		slog.Error("failed to Decode user input", "err", err)
		ResponseWithErr(w, http.StatusInternalServerError, "error on decoding input")
		return
	}
	// later can set up more advanced validation, for not not important
	if userReqeust.FirstName == "" || userReqeust.LastName == "" || userReqeust.Email == "" || userReqeust.Password == "" {
		ResponseWithErr(w, http.StatusBadRequest, "fields cannot be empty")
		return
	}

	//hash pasword
	hash, err := HashPassword(userReqeust.Password)
	if err != nil {
		slog.Error("failed to hash password", "err", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user := users.User{
		ID:        uuid.New(),
		FirstName: userReqeust.FirstName,
		LastName:  userReqeust.LastName,
		Email:     userReqeust.Email,
		Password:  hash,
	}

	err = s.userSrv.Create(r.Context(), user)
	if err != nil {
		slog.Error("failed to create user", "err", err)
		ResponseWithErr(w, 500, "failed to create user")
		return
	}
	ResponseWithJSON(w, http.StatusCreated, nil)
}

func (s *HTTPServer) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, ok := parseUUID(w, idParam)
	if !ok {
		return
	}

	user, err := s.userSrv.Get(r.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			ResponseWithErr(w, http.StatusNotFound, "user not found")
			return
		}
		slog.Error("failed to get user", "err", err)
		return
	}

	userDto := UserDTO{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
	}

	ResponseWithJSON(w, 200, userDto)

}

func (s *HTTPServer) GetMeHandler(w http.ResponseWriter, r *http.Request) {
	idParm := getUserID(r)
	userID, ok := parseUUID(w, idParm)
	if !ok {
		return
	}
	user, err := s.userSrv.Get(r.Context(), userID)
	if err != nil {
		if isNotFoundError(err) {
			ResponseWithErr(w, http.StatusNotFound, "user not found")
			return
		}
		slog.Error("failed to get user", "err", err)
		return
	}

	userDto := UserDTO{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
	}

	ResponseWithJSON(w, 200, userDto)

}

func (s *HTTPServer) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, ok := parseUUID(w, idParam)
	if !ok {
		return
	}

	err := s.userSrv.Delete(r.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			ResponseWithErr(w, http.StatusNotFound, "user not found")
			return
		}
		slog.Error("failed to delete user", "err", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "user deleted successfully"})
}

func (s *HTTPServer) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := s.userSrv.List(r.Context())
	if err != nil {
		slog.Error("failed to list users", "err", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	// Convert users to DTOs to exclude sensitive data
	userDTOs := make([]UserDTO, 0, len(users))
	for _, user := range users {
		userDTOs = append(userDTOs, UserDTO{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			IsAdmin:   user.IsAdmin,
		})
	}

	ResponseWithJSON(w, http.StatusOK, userDTOs)
}
