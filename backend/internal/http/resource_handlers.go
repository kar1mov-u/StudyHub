package http

import (
	"StudyHub/backend/internal/resources"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// get the week from the urlParam, userID from the r.Context().  POST /resources/week/week_id
func (s *HTTPServer) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	weekIDParam := chi.URLParam(r, "week_id")
	weekID, ok := parseUUID(w, weekIDParam)
	if !ok {
		return
	}
	idStr := getUserID(r)
	userID, ok := parseUUID(w, idStr)
	if !ok {
		return
	}
	// r.Body = http.MaxBytesReader(w, r.Body, 100<<20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		ResponseWithErr(w, http.StatusBadRequest, "cannot access form file data")
		log.Println(err)
		return
	}
	defer file.Close()
	err = s.resourceSrv.UploadResource(r.Context(), file, handler.Size, resources.Resource{ID: uuid.New(), WeekID: weekID, UserID: userID, ResourceType: resources.ResourceFile, Name: handler.Filename})
	if err != nil {
		slog.Error("failed to upload file", "err", err)
		ResponseWithErr(w, 500, "failed to upload file")
		return
	}
	ResponseWithJSON(w, http.StatusCreated, nil)

}

func (s *HTTPServer) ListResourcesForWeekHandler(w http.ResponseWriter, r *http.Request) {
	weekIDParam := chi.URLParam(r, "week_id")
	weekID, ok := parseUUID(w, weekIDParam)
	if !ok {
		return
	}

	resources, err := s.resourceSrv.ListResourcesForWeek(r.Context(), weekID)
	if err != nil {
		slog.Error("failed to list resources for week", "err", err)
		ResponseWithErr(w, 500, "failed to list resources")
		return
	}

	ResponseWithJSON(w, 200, resources)

}
