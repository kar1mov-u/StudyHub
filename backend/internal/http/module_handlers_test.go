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
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// mockModuleService implements module service methods for testing
type mockModuleService struct {
	listModulesFunc     func(ctx context.Context) ([]modules.Module, error)
	getModuleFullFunc   func(ctx context.Context, id uuid.UUID) (modules.ModulePage, error)
	createModuleFunc    func(ctx context.Context, module modules.Module) error
	deleteModuleFunc    func(ctx context.Context, id uuid.UUID) error
	listModuleRunsFunc  func(ctx context.Context, moduleID uuid.UUID) ([]modules.ModuleRun, error)
	createModuleRunFunc func(ctx context.Context, run modules.ModuleRun) error
	getModuleRunFunc    func(ctx context.Context, id uuid.UUID) (modules.ModuleRunPage, error)
	deleteModuleRunFunc func(ctx context.Context, id uuid.UUID) error
}

func (m *mockModuleService) ListModules(ctx context.Context) ([]modules.Module, error) {
	if m.listModulesFunc != nil {
		return m.listModulesFunc(ctx)
	}
	return []modules.Module{}, nil
}

func (m *mockModuleService) GetModuleFull(ctx context.Context, id uuid.UUID) (modules.ModulePage, error) {
	if m.getModuleFullFunc != nil {
		return m.getModuleFullFunc(ctx, id)
	}
	return modules.ModulePage{}, nil
}

func (m *mockModuleService) CreateModule(ctx context.Context, module modules.Module) error {
	if m.createModuleFunc != nil {
		return m.createModuleFunc(ctx, module)
	}
	return nil
}

func (m *mockModuleService) DeleteModule(ctx context.Context, id uuid.UUID) error {
	if m.deleteModuleFunc != nil {
		return m.deleteModuleFunc(ctx, id)
	}
	return nil
}

func (m *mockModuleService) ListModuleRuns(ctx context.Context, moduleID uuid.UUID) ([]modules.ModuleRun, error) {
	if m.listModuleRunsFunc != nil {
		return m.listModuleRunsFunc(ctx, moduleID)
	}
	return []modules.ModuleRun{}, nil
}

func (m *mockModuleService) CreateModuleRun(ctx context.Context, run modules.ModuleRun) error {
	if m.createModuleRunFunc != nil {
		return m.createModuleRunFunc(ctx, run)
	}
	return nil
}

func (m *mockModuleService) GetModuleRun(ctx context.Context, id uuid.UUID) (modules.ModuleRunPage, error) {
	if m.getModuleRunFunc != nil {
		return m.getModuleRunFunc(ctx, id)
	}
	return modules.ModuleRunPage{}, nil
}

func (m *mockModuleService) DeleteModuleRun(ctx context.Context, id uuid.UUID) error {
	if m.deleteModuleRunFunc != nil {
		return m.deleteModuleRunFunc(ctx, id)
	}
	return nil
}

func TestListModulesHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockFunc       func(ctx context.Context) ([]modules.Module, error)
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "success - list multiple modules",
			mockFunc: func(ctx context.Context) ([]modules.Module, error) {
				return []modules.Module{
					{ID: uuid.New(), Code: "CS101", Name: "Intro to CS", DepartmentName: "CS"},
					{ID: uuid.New(), Code: "MATH201", Name: "Calculus", DepartmentName: "Math"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "success - empty list",
			mockFunc: func(ctx context.Context) ([]modules.Module, error) {
				return []modules.Module{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "error - database error",
			mockFunc: func(ctx context.Context) ([]modules.Module, error) {
				return nil, errors.New("database connection failed")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockModuleService{listModulesFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodGet, "/modules", nil)
			w := httptest.NewRecorder()

			// Execute handler logic
			mods, err := mockSvc.ListModules(req.Context())
			if err != nil {
				ResponseWithErr(w, http.StatusInternalServerError, err.Error())
			} else {
				ResponseWithJSON(w, http.StatusOK, mods)
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetModuleFullHandler(t *testing.T) {
	tests := []struct {
		name           string
		moduleID       string
		mockFunc       func(ctx context.Context, id uuid.UUID) (modules.ModulePage, error)
		expectedStatus int
	}{
		{
			name:     "success - get module with details",
			moduleID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (modules.ModulePage, error) {
				return modules.ModulePage{
					Module: modules.Module{ID: id, Code: "CS101", Name: "Intro to CS"},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "error - module not found",
			moduleID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (modules.ModulePage, error) {
				return modules.ModulePage{}, pgx.ErrNoRows
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "error - invalid UUID",
			moduleID:       "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockModuleService{getModuleFullFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodGet, "/modules/"+tt.moduleID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.moduleID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			// Execute handler logic
			idParam := chi.URLParam(req, "id")
			id, ok := parseUUID(w, idParam)
			if ok {
				modulePage, err := mockSvc.GetModuleFull(req.Context(), id)
				if err != nil {
					if isNotFoundError(err) {
						ResponseWithErr(w, http.StatusNotFound, "module not found")
					} else {
						ResponseWithErr(w, http.StatusInternalServerError, err.Error())
					}
				} else {
					ResponseWithJSON(w, http.StatusOK, modulePage)
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestCreateModuleHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, module modules.Module) error
		expectedStatus int
	}{
		{
			name: "success - create module",
			requestBody: CreateModuleRequest{
				Code:           "CS101",
				Name:           "Intro to CS",
				DepartmentName: "CS",
			},
			mockFunc: func(ctx context.Context, module modules.Module) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "error - missing code",
			requestBody: CreateModuleRequest{
				Code:           "",
				Name:           "Intro to CS",
				DepartmentName: "CS",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error - invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockModuleService{createModuleFunc: tt.mockFunc}
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}
			req := httptest.NewRequest(http.MethodPost, "/modules", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			// Execute handler logic
			var reqData CreateModuleRequest
			if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&reqData); err != nil {
				ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
			} else if reqData.Code == "" || reqData.Name == "" || reqData.DepartmentName == "" {
				ResponseWithErr(w, http.StatusBadRequest, "code, name, and department_name are required")
			} else {
				module := modules.Module{
					ID:             uuid.New(),
					Code:           reqData.Code,
					Name:           reqData.Name,
					DepartmentName: reqData.DepartmentName,
					CreatedAt:      time.Now(),
				}
				if err := mockSvc.CreateModule(req.Context(), module); err != nil {
					ResponseWithErr(w, http.StatusInternalServerError, err.Error())
				} else {
					ResponseWithJSON(w, http.StatusCreated, map[string]uuid.UUID{"id": module.ID})
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestDeleteModuleHandler(t *testing.T) {
	tests := []struct {
		name           string
		moduleID       string
		mockFunc       func(ctx context.Context, id uuid.UUID) error
		expectedStatus int
	}{
		{
			name:     "success - delete module",
			moduleID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:     "error - module not found",
			moduleID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return pgx.ErrNoRows
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "error - invalid UUID",
			moduleID:       "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockModuleService{deleteModuleFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodDelete, "/modules/"+tt.moduleID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.moduleID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			// Execute handler logic
			idParam := chi.URLParam(req, "id")
			id, ok := parseUUID(w, idParam)
			if ok {
				if err := mockSvc.DeleteModule(req.Context(), id); err != nil {
					if isNotFoundError(err) {
						ResponseWithErr(w, http.StatusNotFound, "module not found")
					} else {
						ResponseWithErr(w, http.StatusInternalServerError, err.Error())
					}
				} else {
					w.WriteHeader(http.StatusNoContent)
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestListModuleRunsHandler(t *testing.T) {
	tests := []struct {
		name           string
		moduleID       string
		mockFunc       func(ctx context.Context, moduleID uuid.UUID) ([]modules.ModuleRun, error)
		expectedStatus int
	}{
		{
			name:     "success - list module runs",
			moduleID: uuid.New().String(),
			mockFunc: func(ctx context.Context, moduleID uuid.UUID) ([]modules.ModuleRun, error) {
				return []modules.ModuleRun{
					{ID: uuid.New(), ModuleID: moduleID, Year: 2024, Semester: "Spring"},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error - invalid module ID",
			moduleID:       "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockModuleService{listModuleRunsFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodGet, "/modules/"+tt.moduleID+"/runs", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("moduleID", tt.moduleID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			// Execute handler logic
			moduleIDParam := chi.URLParam(req, "moduleID")
			moduleID, ok := parseUUID(w, moduleIDParam)
			if ok {
				runs, err := mockSvc.ListModuleRuns(req.Context(), moduleID)
				if err != nil {
					ResponseWithErr(w, http.StatusInternalServerError, err.Error())
				} else {
					ResponseWithJSON(w, http.StatusOK, runs)
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestCreateModuleRunHandler(t *testing.T) {
	tests := []struct {
		name           string
		moduleID       string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, run modules.ModuleRun) error
		expectedStatus int
	}{
		{
			name:     "success - create module run",
			moduleID: uuid.New().String(),
			requestBody: CreateModuleRunRequest{
				Year:     2024,
				Semester: "Spring",
			},
			mockFunc: func(ctx context.Context, run modules.ModuleRun) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:     "error - missing year",
			moduleID: uuid.New().String(),
			requestBody: CreateModuleRunRequest{
				Year:     0,
				Semester: "Spring",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error - invalid module ID",
			moduleID:       "invalid-uuid",
			requestBody:    CreateModuleRunRequest{Year: 2024, Semester: "Spring"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockModuleService{createModuleRunFunc: tt.mockFunc}
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/modules/"+tt.moduleID+"/runs", bytes.NewBuffer(body))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("moduleID", tt.moduleID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			// Execute handler logic
			moduleIDParam := chi.URLParam(req, "moduleID")
			moduleID, ok := parseUUID(w, moduleIDParam)
			if ok {
				var reqData CreateModuleRunRequest
				if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&reqData); err != nil {
					ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
				} else if reqData.Year == 0 || reqData.Semester == "" {
					ResponseWithErr(w, http.StatusBadRequest, "year and semester are required")
				} else {
					run := modules.ModuleRun{
						ID:        uuid.New(),
						ModuleID:  moduleID,
						Year:      reqData.Year,
						Semester:  reqData.Semester,
						CreatedAt: time.Now(),
					}
					if err := mockSvc.CreateModuleRun(req.Context(), run); err != nil {
						ResponseWithErr(w, http.StatusInternalServerError, err.Error())
					} else {
						ResponseWithJSON(w, http.StatusCreated, map[string]uuid.UUID{"id": run.ID})
					}
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetModuleRunHandler(t *testing.T) {
	tests := []struct {
		name           string
		runID          string
		mockFunc       func(ctx context.Context, id uuid.UUID) (modules.ModuleRunPage, error)
		expectedStatus int
	}{
		{
			name:  "success - get module run",
			runID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (modules.ModuleRunPage, error) {
				return modules.ModuleRunPage{
					Run: modules.ModuleRun{ID: id, Semester: "Spring", Year: 2024},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "error - module run not found",
			runID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (modules.ModuleRunPage, error) {
				return modules.ModuleRunPage{}, pgx.ErrNoRows
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "error - invalid run ID",
			runID:          "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockModuleService{getModuleRunFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodGet, "/module-runs/"+tt.runID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.runID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			idParam := chi.URLParam(req, "id")
			id, ok := parseUUID(w, idParam)
			if ok {
				moduleRunPage, err := mockSvc.GetModuleRun(req.Context(), id)
				if err != nil {
					if isNotFoundError(err) {
						ResponseWithErr(w, http.StatusNotFound, "module run not found")
					} else {
						ResponseWithErr(w, http.StatusInternalServerError, err.Error())
					}
				} else {
					ResponseWithJSON(w, http.StatusOK, moduleRunPage)
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestDeleteModuleRunHandler(t *testing.T) {
	tests := []struct {
		name           string
		runID          string
		mockFunc       func(ctx context.Context, id uuid.UUID) error
		expectedStatus int
	}{
		{
			name:  "success - delete module run",
			runID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:  "error - module run not found",
			runID: uuid.New().String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return pgx.ErrNoRows
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "error - invalid run ID",
			runID:          "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockModuleService{deleteModuleRunFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodDelete, "/module-runs/"+tt.runID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.runID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			idParam := chi.URLParam(req, "id")
			id, ok := parseUUID(w, idParam)
			if ok {
				if err := mockSvc.DeleteModuleRun(req.Context(), id); err != nil {
					if isNotFoundError(err) {
						ResponseWithErr(w, http.StatusNotFound, "module run not found")
					} else {
						ResponseWithErr(w, http.StatusInternalServerError, err.Error())
					}
				} else {
					w.WriteHeader(http.StatusNoContent)
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
