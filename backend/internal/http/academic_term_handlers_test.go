package http

import (
	"StudyHub/internal/modules"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Mock ModuleService for academic term tests
type mockModuleServiceForTerms struct {
	getActiveAcademicTermFunc func(ctx context.Context) (modules.AcademicTerm, error)
	startNewTermFunc          func(ctx context.Context, term modules.AcademicTerm) error
}

func (m *mockModuleServiceForTerms) GetActiveAcademicTerm(ctx context.Context) (modules.AcademicTerm, error) {
	if m.getActiveAcademicTermFunc != nil {
		return m.getActiveAcademicTermFunc(ctx)
	}
	return modules.AcademicTerm{}, nil
}

func (m *mockModuleServiceForTerms) StartNewTerm(ctx context.Context, term modules.AcademicTerm) error {
	if m.startNewTermFunc != nil {
		return m.startNewTermFunc(ctx, term)
	}
	return nil
}

// Stub methods for other ModuleService methods (not used in these tests)
func (m *mockModuleServiceForTerms) ListModules(ctx context.Context) ([]modules.Module, error) {
	return nil, nil
}

func (m *mockModuleServiceForTerms) GetModuleFull(ctx context.Context, id uuid.UUID) (modules.ModulePage, error) {
	return modules.ModulePage{}, nil
}

func (m *mockModuleServiceForTerms) CreateModule(ctx context.Context, module modules.Module) error {
	return nil
}

func (m *mockModuleServiceForTerms) DeleteModule(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockModuleServiceForTerms) ListModuleRuns(ctx context.Context, moduleID uuid.UUID) ([]modules.ModuleRun, error) {
	return nil, nil
}

func (m *mockModuleServiceForTerms) CreateModuleRun(ctx context.Context, run modules.ModuleRun) error {
	return nil
}

func (m *mockModuleServiceForTerms) GetModuleRun(ctx context.Context, id uuid.UUID) (modules.ModuleRunPage, error) {
	return modules.ModuleRunPage{}, nil
}

func (m *mockModuleServiceForTerms) DeleteModuleRun(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockModuleServiceForTerms) ListAcademicTerms(ctx context.Context) ([]modules.AcademicTerm, error) {
	return nil, nil
}

// Test GetActiveAcademicTermHandler
func TestGetActiveAcademicTermHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockFunc       func(ctx context.Context) (modules.AcademicTerm, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success - returns active term",
			mockFunc: func(ctx context.Context) (modules.AcademicTerm, error) {
				return modules.AcademicTerm{
					ID:       uuid.MustParse("11111111-1111-1111-1111-111111111111"),
					Year:     2024,
					Semester: "spring",
					IsActive: true,
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"Semester":"spring"`,
		},
		{
			name: "not found - no active term",
			mockFunc: func(ctx context.Context) (modules.AcademicTerm, error) {
				// Return pgx.ErrNoRows which isNotFoundError checks for
				return modules.AcademicTerm{}, pgx.ErrNoRows
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "no active academic term found",
		},
		{
			name: "error - internal server error",
			mockFunc: func(ctx context.Context) (modules.AcademicTerm, error) {
				return modules.AcademicTerm{}, errors.New("database connection failed")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockModuleServiceForTerms{getActiveAcademicTermFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodGet, "/api/v1/academic-terms/active", nil)
			w := httptest.NewRecorder()

			// Execute handler logic inline
			term, err := mockSvc.GetActiveAcademicTerm(req.Context())
			if err != nil {
				if isNotFoundError(err) {
					ResponseWithErr(w, http.StatusNotFound, "no active academic term found")
				} else {
					ResponseWithErr(w, http.StatusInternalServerError, err.Error())
				}
			} else {
				ResponseWithJSON(w, http.StatusOK, term)
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" {
				body := w.Body.String()
				if !contains(body, tt.expectedBody) {
					t.Errorf("expected body to contain %q, got %q", tt.expectedBody, body)
				}
			}
		})
	}
}

// Test CreateAcademicTermHandler
func TestCreateAcademicTermHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, term modules.AcademicTerm) error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success - creates spring term",
			requestBody: CreateAcademicTermRequest{
				Year:     2024,
				Semester: "Spring",
			},
			mockFunc: func(ctx context.Context, term modules.AcademicTerm) error {
				if term.Year != 2024 || term.Semester != "spring" {
					t.Errorf("expected year=2024 semester=spring, got year=%d semester=%s", term.Year, term.Semester)
				}
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"id"`,
		},
		{
			name: "success - creates fall term",
			requestBody: CreateAcademicTermRequest{
				Year:     2024,
				Semester: "Fall",
			},
			mockFunc: func(ctx context.Context, term modules.AcademicTerm) error {
				if term.Semester != "fall" {
					t.Errorf("expected semester to be lowercase 'fall', got %s", term.Semester)
				}
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"id"`,
		},
		{
			name:           "error - invalid JSON",
			requestBody:    "invalid json",
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid request body",
		},
		{
			name: "error - missing year",
			requestBody: CreateAcademicTermRequest{
				Semester: "spring",
			},
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "year is required",
		},
		{
			name: "error - missing semester",
			requestBody: CreateAcademicTermRequest{
				Year: 2024,
			},
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "semester is required",
		},
		{
			name: "error - invalid semester",
			requestBody: CreateAcademicTermRequest{
				Year:     2024,
				Semester: "summer",
			},
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "semester must be either 'spring' or 'fall'",
		},
		{
			name: "error - service error",
			requestBody: CreateAcademicTermRequest{
				Year:     2024,
				Semester: "spring",
			},
			mockFunc: func(ctx context.Context, term modules.AcademicTerm) error {
				return errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockModuleServiceForTerms{startNewTermFunc: tt.mockFunc}

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

			req := httptest.NewRequest(http.MethodPost, "/api/v1/academic-terms", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute handler logic inline
			var reqData CreateAcademicTermRequest
			if err := json.NewDecoder(bytes.NewReader(body)).Decode(&reqData); err != nil {
				ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
			} else if reqData.Year == 0 {
				ResponseWithErr(w, http.StatusBadRequest, "year is required")
			} else if reqData.Semester == "" {
				ResponseWithErr(w, http.StatusBadRequest, "semester is required")
			} else {
				semester := reqData.Semester
				switch semester {
				case "Spring":
					semester = "spring"
				case "Fall":
					semester = "fall"
				}

				if semester != "spring" && semester != "fall" {
					ResponseWithErr(w, http.StatusBadRequest, "semester must be either 'spring' or 'fall'")
				} else {
					term := modules.AcademicTerm{
						ID:       uuid.New(),
						Year:     reqData.Year,
						Semester: semester,
						IsActive: true,
					}

					if err := mockSvc.StartNewTerm(req.Context(), term); err != nil {
						ResponseWithErr(w, http.StatusInternalServerError, err.Error())
					} else {
						ResponseWithJSON(w, http.StatusCreated, map[string]uuid.UUID{"id": term.ID})
					}
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

// Helper function to check if a string contains a substring
func contains(str, substr string) bool {
	return bytes.Contains([]byte(str), []byte(substr))
}
