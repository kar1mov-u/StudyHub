package http

import (
	"StudyHub/backend/internal/modules"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// Request DTOs
type CreateAcademicTermRequest struct {
	Year     int    `json:"year"`
	Semester string `json:"semester"`
}

// Handler 1: Get active academic term
func (s *HTTPServer) GetActiveAcademicTermHandler(w http.ResponseWriter, r *http.Request) {
	term, err := s.moduleSrv.GetActiveAcademicTerm(r.Context())
	if err != nil {
		if isNotFoundError(err) {
			ResponseWithErr(w, http.StatusNotFound, "no active academic term found")
			return
		}
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseWithJSON(w, http.StatusOK, term)
}

// Handler 2: List all academic terms
func (s *HTTPServer) ListAcademicTermsHandler(w http.ResponseWriter, r *http.Request) {
	terms, err := s.moduleSrv.ListAcademicTerms(r.Context())
	if err != nil {
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseWithJSON(w, http.StatusOK, terms)
}

// Handler 3: Create a new academic term
func (s *HTTPServer) CreateAcademicTermHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateAcademicTermRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validation
	if req.Year == 0 {
		ResponseWithErr(w, http.StatusBadRequest, "year is required")
		return
	}

	if req.Semester == "" {
		ResponseWithErr(w, http.StatusBadRequest, "semester is required")
		return
	}

	// Validate semester - only "spring" or "fall" allowed
	semester := strings.ToLower(req.Semester)
	if semester != "spring" && semester != "fall" {
		ResponseWithErr(w, http.StatusBadRequest, "semester must be either 'spring' or 'fall'")
		return
	}

	term := modules.AcademicTerm{
		ID:       uuid.New(),
		Year:     req.Year,
		Semester: semester,
		IsActive: false,
	}

	if err := s.moduleSrv.CreateAcademicTerm(r.Context(), term); err != nil {
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseWithJSON(w, http.StatusCreated, map[string]uuid.UUID{"id": term.ID})
}

// Handler 4: Deactivate an academic term
func (s *HTTPServer) DeactivateAcademicTermHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, ok := parseUUID(w, idParam)
	if !ok {
		return
	}

	if err := s.moduleSrv.DeactivateAcademicTerm(r.Context(), id); err != nil {
		if isNotFoundError(err) {
			ResponseWithErr(w, http.StatusNotFound, "academic term not found")
			return
		}
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *HTTPServer) ActivateAcademicTermHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, ok := parseUUID(w, idParam)
	if !ok {
		return
	}

	if err := s.moduleSrv.ActivateAcademicTerm(r.Context(), id); err != nil {
		if isNotFoundError(err) {
			ResponseWithErr(w, http.StatusNotFound, "academic term not found")
			return
		}
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
