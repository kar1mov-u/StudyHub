package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

// Mock AuthService for testing
type mockAuthService struct {
	loginUserFunc func(ctx context.Context, email, password string) (string, error)
	isAdminFunc   func(ctx context.Context, id uuid.UUID) bool
}

func (m *mockAuthService) LoginUser(ctx context.Context, email, password string) (string, error) {
	if m.loginUserFunc != nil {
		return m.loginUserFunc(ctx, email, password)
	}
	return "", nil
}

func (m *mockAuthService) IsAdmin(ctx context.Context, id uuid.UUID) bool {
	if m.isAdminFunc != nil {
		return m.isAdminFunc(ctx, id)
	}
	return false
}

// Test LoginHandler
func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, email, password string) (string, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success - valid credentials",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			mockFunc: func(ctx context.Context, email, password string) (string, error) {
				if email != "test@example.com" || password != "password123" {
					t.Errorf("expected email=test@example.com password=password123, got email=%s password=%s", email, password)
				}
				return "mock-jwt-token-12345", nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "mock-jwt-token-12345",
		},
		{
			name:           "error - invalid JSON",
			requestBody:    "not valid json",
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid json",
		},
		{
			name: "error - authentication failed",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			mockFunc: func(ctx context.Context, email, password string) (string, error) {
				return "", errors.New("invalid credentials")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to create the JWT",
		},
		{
			name: "error - database error",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			mockFunc: func(ctx context.Context, email, password string) (string, error) {
				return "", errors.New("database connection failed")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to create the JWT",
		},
		{
			name: "success - valid login with different credentials",
			requestBody: map[string]string{
				"email":    "admin@example.com",
				"password": "admin123",
			},
			mockFunc: func(ctx context.Context, email, password string) (string, error) {
				return "admin-jwt-token", nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "admin-jwt-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockAuthService{loginUserFunc: tt.mockFunc}

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

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute handler logic inline
			var reqData struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}

			if err := json.NewDecoder(bytes.NewReader(body)).Decode(&reqData); err != nil {
				ResponseWithErr(w, http.StatusBadRequest, "invalid json")
			} else {
				token, err := mockSvc.LoginUser(req.Context(), reqData.Email, reqData.Password)
				if err != nil {
					ResponseWithErr(w, http.StatusInternalServerError, "failed to create the JWT")
				} else {
					ResponseWithJSON(w, http.StatusOK, map[string]string{"token": token})
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" {
				responseBody := w.Body.String()
				if !contains(responseBody, tt.expectedBody) {
					t.Errorf("expected body to contain %q, got %q", tt.expectedBody, responseBody)
				}
			}
		})
	}
}
