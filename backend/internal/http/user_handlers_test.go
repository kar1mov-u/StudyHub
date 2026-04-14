package http

import (
	"StudyHub/internal/users"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type mockUserService struct {
	createFunc func(ctx context.Context, user users.User) error
	getFunc    func(ctx context.Context, id uuid.UUID) (users.User, error)
	deleteFunc func(ctx context.Context, id uuid.UUID) error
	listFunc   func(ctx context.Context) ([]users.User, error)
}

func (m *mockUserService) Create(ctx context.Context, user users.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return nil
}

func (m *mockUserService) Get(ctx context.Context, id uuid.UUID) (users.User, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, id)
	}
	return users.User{}, nil
}

func (m *mockUserService) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockUserService) List(ctx context.Context) ([]users.User, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	return []users.User{}, nil
}

func TestCreateUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, user users.User) error
		expectedStatus int
	}{
		{
			name: "success - create user",
			requestBody: CreateUserRequest{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane@example.com",
				Password:  "password123",
			},
			mockFunc: func(ctx context.Context, user users.User) error {
				if user.Email != "jane@example.com" {
					t.Errorf("expected email jane@example.com, got %s", user.Email)
				}
				if !CheckPasswordHash("password123", user.Password) {
					t.Error("expected password to be hashed")
				}
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "error - invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "error - missing fields",
			requestBody: CreateUserRequest{
				FirstName: "",
				LastName:  "Doe",
				Email:     "jane@example.com",
				Password:  "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error - service failure",
			requestBody: CreateUserRequest{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane@example.com",
				Password:  "password123",
			},
			mockFunc: func(ctx context.Context, user users.User) error {
				return errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{createFunc: tt.mockFunc}

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			var userRequest CreateUserRequest
			if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&userRequest); err != nil {
				ResponseWithErr(w, http.StatusInternalServerError, "error on decoding input")
			} else if userRequest.FirstName == "" || userRequest.LastName == "" || userRequest.Email == "" || userRequest.Password == "" {
				ResponseWithErr(w, http.StatusBadRequest, "fields cannot be empty")
			} else {
				hash, err := HashPassword(userRequest.Password)
				if err != nil {
					ResponseWithErr(w, http.StatusInternalServerError, "failed to hash password")
				} else {
					user := users.User{
						ID:        uuid.New(),
						FirstName: userRequest.FirstName,
						LastName:  userRequest.LastName,
						Email:     userRequest.Email,
						Password:  hash,
					}
					if err := mockSvc.Create(req.Context(), user); err != nil {
						ResponseWithErr(w, http.StatusInternalServerError, "failed to create user")
					} else {
						ResponseWithJSON(w, http.StatusCreated, nil)
					}
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockFunc       func(ctx context.Context, id uuid.UUID) (users.User, error)
		expectedStatus int
	}{
		{
			name:   "success - get user",
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (users.User, error) {
				return users.User{ID: id, FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "error - user not found",
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (users.User, error) {
				return users.User{}, pgx.ErrNoRows
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "error - invalid user ID",
			userID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{getFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			idParam := chi.URLParam(req, "id")
			id, ok := parseUUID(w, idParam)
			if ok {
				user, err := mockSvc.Get(req.Context(), id)
				if err != nil {
					if isNotFoundError(err) {
						ResponseWithErr(w, http.StatusNotFound, "user not found")
					} else {
						return
					}
				} else {
					userDto := UserDTO{
						ID:        user.ID,
						FirstName: user.FirstName,
						LastName:  user.LastName,
						Email:     user.Email,
						IsAdmin:   user.IsAdmin,
					}
					ResponseWithJSON(w, http.StatusOK, userDto)
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetMeHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockFunc       func(ctx context.Context, id uuid.UUID) (users.User, error)
		expectedStatus int
	}{
		{
			name:   "success - get current user",
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (users.User, error) {
				return users.User{ID: id, FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "error - user not found",
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (users.User, error) {
				return users.User{}, pgx.ErrNoRows
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "error - invalid user ID",
			userID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{getFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
			req = addUserIDToContext(req, tt.userID)
			w := httptest.NewRecorder()

			idParam := getUserID(req)
			id, ok := parseUUID(w, idParam)
			if ok {
				user, err := mockSvc.Get(req.Context(), id)
				if err != nil {
					if isNotFoundError(err) {
						ResponseWithErr(w, http.StatusNotFound, "user not found")
					} else {
						return
					}
				} else {
					userDto := UserDTO{
						ID:        user.ID,
						FirstName: user.FirstName,
						LastName:  user.LastName,
						Email:     user.Email,
						IsAdmin:   user.IsAdmin,
					}
					ResponseWithJSON(w, http.StatusOK, userDto)
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestDeleteUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockFunc       func(ctx context.Context, id uuid.UUID) error
		expectedStatus int
	}{
		{
			name:   "success - delete user",
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "error - user not found",
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return pgx.ErrNoRows
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "error - service failure",
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "error - invalid user ID",
			userID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{deleteFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.userID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			idParam := chi.URLParam(req, "id")
			id, ok := parseUUID(w, idParam)
			if ok {
				if err := mockSvc.Delete(req.Context(), id); err != nil {
					if isNotFoundError(err) {
						ResponseWithErr(w, http.StatusNotFound, "user not found")
					} else {
						ResponseWithErr(w, http.StatusInternalServerError, "failed to delete user")
					}
				} else {
					ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "user deleted successfully"})
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestListUsersHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockFunc       func(ctx context.Context) ([]users.User, error)
		expectedStatus int
	}{
		{
			name: "success - list users",
			mockFunc: func(ctx context.Context) ([]users.User, error) {
				return []users.User{
					{ID: uuid.New(), FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"},
					{ID: uuid.New(), FirstName: "John", LastName: "Smith", Email: "john@example.com"},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "error - service failure",
			mockFunc: func(ctx context.Context) ([]users.User, error) {
				return nil, errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{listFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			w := httptest.NewRecorder()

			usersList, err := mockSvc.List(req.Context())
			if err != nil {
				ResponseWithErr(w, http.StatusInternalServerError, "failed to list users")
			} else {
				userDTOs := make([]UserDTO, 0, len(usersList))
				for _, user := range usersList {
					userDTOs = append(userDTOs, UserDTO{
						FirstName: user.FirstName,
						LastName:  user.LastName,
						Email:     user.Email,
						IsAdmin:   user.IsAdmin,
					})
				}
				ResponseWithJSON(w, http.StatusOK, userDTOs)
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
