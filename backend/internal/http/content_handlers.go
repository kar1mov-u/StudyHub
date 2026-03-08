package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type ListCardsForObjectsRequest struct {
	IDs []string `json:"ids"`
}

func (s *HTTPServer) ListCardsFromObjects(w http.ResponseWriter, r *http.Request) {
	var payload ListCardsForObjectsRequest

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		slog.Error("failed to Decode Json", "error", err)
		ResponseWithErr(w, http.StatusBadRequest, "invalid request")
		return
	}

	uuids := make([]uuid.UUID, len(payload.IDs))
	for i := range payload.IDs {
		id, err := uuid.Parse(payload.IDs[i])
		if err != nil {
			slog.Error("failed to parse id", "error", err)
			ResponseWithErr(w, http.StatusBadRequest, "invalid objectID")
			return
		}
		uuids = append(uuids, id)
	}

	cards, err := s.contentSrv.ListSelectedObjectsCards(r.Context(), uuids)
	if err != nil {
		slog.Error("failed to list cards", "error", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to list cards")
		return
	}
	ResponseWithJSON(w, 200, cards)
}
