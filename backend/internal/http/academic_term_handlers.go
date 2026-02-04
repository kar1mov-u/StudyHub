package http

import (
	"StudyHub/backend/internal/modules"
	"encoding/json"
	"net/http"
	"strings"

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
		IsActive: true,
	}

	if err := s.moduleSrv.StartNewTerm(r.Context(), term); err != nil {
		ResponseWithErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseWithJSON(w, http.StatusCreated, map[string]uuid.UUID{"id": term.ID})
}
