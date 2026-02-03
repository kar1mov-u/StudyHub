package http

import (
	"StudyHub/backend/internal/resources"
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// get the week from the urlParam, userID from the r.Context().  POST /resources/week/week_id
func (s *HTTPServer) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting upload")
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
		if errors.Is(err, resources.ErrResourceExists) {
			ResponseWithErr(w, http.StatusBadRequest, "file is uploaded by other user already")
			return
		}
		slog.Error("failed to upload file", "err", err)
		ResponseWithErr(w, 500, "failed to upload file")
		return
	}
	ResponseWithJSON(w, http.StatusCreated, nil)
}

func (s *HTTPServer) CreateLinkResource(w http.ResponseWriter, r *http.Request) {
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
	var request CreateLinkResourceRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		slog.Error("failed to decode JSON", "err", err)
		return
	}

	resource := resources.Resource{ID: uuid.New(), WeekID: weekID, UserID: userID, ResourceType: resources.ResourceLink, ExternalLink: &request.URL, Name: request.Name}
	err = s.resourceSrv.CreateLinkResource(r.Context(), resource)
	if err != nil {
		if errors.Is(err, resources.ErrResourceExists) {
			ResponseWithErr(w, http.StatusBadRequest, "Link is already uploaded by other user")
			return
		} else {
			ResponseWithErr(w, http.StatusInternalServerError, "failed to create the Link Reosurce")
			return
		}
	}
	ResponseWithJSON(w, 201, nil)

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
func (s *HTTPServer) ListResourcesForUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "user_id")
	userID, ok := parseUUID(w, userIDParam)
	if !ok {
		return
	}

	resources, err := s.resourceSrv.ListResourceForUser(r.Context(), userID)
	if err != nil {
		slog.Error("failed to list resources for user", "err", err)
		ResponseWithErr(w, 500, "failed to list resources")
		return
	}

	ResponseWithJSON(w, 200, resources)
}

func (s *HTTPServer) GetResourceHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")
	resourceID, ok := parseUUID(w, idParam)
	if !ok {
		return
	}

	url, err := s.resourceSrv.GetResource(r.Context(), resourceID)
	if err != nil {
		ResponseWithErr(w, 500, "failed to get resource")
		return
	}

	ResponseWithJSON(w, 200, map[string]string{"url": url})
}

func (s *HTTPServer) CleanOrphanObjectsHandler(w http.ResponseWriter, r *http.Request) {
	//here try to do some authorization maybe by some token key
	ids, err := s.resourceSrv.CleanOrphanObjects(r.Context())
	if err != nil {
		slog.Error("failed to clean orphan objects", "err", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to delete ")
		return
	}
	slog.Info("cleaned orphan objects from the storage", "id's:", ids)
	ResponseWithJSON(w, 200, nil)
}
func (s *HTTPServer) DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	resourceIDParm := chi.URLParam(r, "id")
	resourceID, ok := parseUUID(w, resourceIDParm)
	if !ok {
		return
	}
	userIDStr := getUserID(r)
	userID, ok := parseUUID(w, userIDStr)
	if !ok {
		return
	}

	err := s.resourceSrv.DeleteResource(context.Background(), userID, resourceID)
	if err != nil {
		slog.Error("failed to delete resource", "err:", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to delete resource")
		return
	}
	ResponseWithJSON(w, 200, nil)
}

type CreateLinkResourceRequest struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}
