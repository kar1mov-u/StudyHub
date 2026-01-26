package http

import (
	"StudyHub/backend/internal/modules"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// Request DTOs
type CreateModuleRequest struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	DepartmentName string `json:"department_name"`
}

type CreateModuleRunRequest struct {
	Year     int    `json:"year"`
	Semester string `json:"semester"`
	IsActive bool   `json:"is_active"`
}

// Helper functions

// Handler 1: List all modules
func (s *HTTPServer) ListModulesHandler(w http.ResponseWriter, r *http.Request) {
	modules, err := s.moduleSrv.ListModules(r.Context())
	if err != nil {
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseWithJSON(w, http.StatusOK, modules)
}

// Handler 2: Get module with full details (module + active run + weeks)
func (s *HTTPServer) GetModuleFullHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, ok := parseUUID(w, idParam)
	if !ok {
		return
	}

	modulePage, err := s.moduleSrv.GetModuleFull(r.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			ResponseWithErr(w, http.StatusNotFound, "module not found")
			return
		}
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseWithJSON(w, http.StatusOK, modulePage)
}

// Handler 3: Create a new module
func (s *HTTPServer) CreateModuleHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Basic validation
	if req.Code == "" || req.Name == "" || req.DepartmentName == "" {
		ResponseWithErr(w, http.StatusBadRequest, "code, name, and department_name are required")
		return
	}

	module := modules.Module{
		ID:             uuid.New(),
		Code:           req.Code,
		Name:           req.Name,
		DepartmentName: req.DepartmentName,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.moduleSrv.CreateModule(r.Context(), module); err != nil {
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseWithJSON(w, http.StatusCreated, map[string]uuid.UUID{"id": module.ID})
}

// Handler 4: Delete a module
func (s *HTTPServer) DeleteModuleHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, ok := parseUUID(w, idParam)
	if !ok {
		return
	}

	if err := s.moduleSrv.DeleteModule(r.Context(), id); err != nil {
		if isNotFoundError(err) {
			ResponseWithErr(w, http.StatusNotFound, "module not found")
			return
		}
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Handler 5: List all runs for a specific module
func (s *HTTPServer) ListModuleRunsHandler(w http.ResponseWriter, r *http.Request) {
	moduleIDParam := chi.URLParam(r, "moduleID")
	moduleID, ok := parseUUID(w, moduleIDParam)
	if !ok {
		return
	}

	runs, err := s.moduleSrv.ListModuleRuns(r.Context(), moduleID)
	if err != nil {
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseWithJSON(w, http.StatusOK, runs)
}

// Handler 6: Create a new module run
func (s *HTTPServer) CreateModuleRunHandler(w http.ResponseWriter, r *http.Request) {
	moduleIDParam := chi.URLParam(r, "moduleID")
	moduleID, ok := parseUUID(w, moduleIDParam)
	if !ok {
		return
	}

	var req CreateModuleRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Basic validation
	if req.Year == 0 || req.Semester == "" {
		ResponseWithErr(w, http.StatusBadRequest, "year and semester are required")
		return
	}

	moduleRun := modules.ModuleRun{
		ID:        uuid.New(),
		ModuleID:  moduleID,
		Year:      req.Year,
		Semester:  req.Semester,
		IsActive:  req.IsActive,
		CreatedAt: time.Now(),
	}

	if err := s.moduleSrv.CreateModuleRun(r.Context(), moduleRun); err != nil {
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseWithJSON(w, http.StatusCreated, map[string]uuid.UUID{"id": moduleRun.ID})
}

// Handler 7: Get module run with weeks
func (s *HTTPServer) GetModuleRunHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, ok := parseUUID(w, idParam)
	if !ok {
		return
	}

	moduleRunPage, err := s.moduleSrv.GetModuleRun(r.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			ResponseWithErr(w, http.StatusNotFound, "module run not found")
			return
		}
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseWithJSON(w, http.StatusOK, moduleRunPage)
}

// Handler 8: Delete a module run
func (s *HTTPServer) DeleteModuleRunHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, ok := parseUUID(w, idParam)
	if !ok {
		return
	}

	if err := s.moduleSrv.DeleteModuleRun(r.Context(), id); err != nil {
		if isNotFoundError(err) {
			ResponseWithErr(w, http.StatusNotFound, "module run not found")
			return
		}
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
