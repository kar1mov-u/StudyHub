package http

import (
	"StudyHub/internal/content"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
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

// AddCardToDeckHandler adds an auto-generated flashcard to user's deck
func (s *HTTPServer) AddCardToDeckHandler(w http.ResponseWriter, r *http.Request) {
	weekIDParam := chi.URLParam(r, "week_id")
	weekID, ok := parseUUID(w, weekIDParam)
	if !ok {
		return
	}

	userIDStr := getUserID(r)
	userID, ok := parseUUID(w, userIDStr)
	if !ok {
		return
	}

	var req content.AddCardToDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", "error", err)
		ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	flashcardID, err := uuid.Parse(req.FlashcardID)
	if err != nil {
		ResponseWithErr(w, http.StatusBadRequest, "invalid flashcard_id")
		return
	}

	err = s.contentSrv.AddCardToUserDeck(r.Context(), userID, weekID, flashcardID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ResponseWithErr(w, http.StatusNotFound, err.Error())
			return
		}
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			ResponseWithErr(w, http.StatusConflict, "card already in deck")
			return
		}
		slog.Error("failed to add card to deck", "error", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to add card to deck")
		return
	}

	ResponseWithJSON(w, http.StatusCreated, nil)
}

// CreateCustomCardHandler creates a custom flashcard in user's deck
func (s *HTTPServer) CreateCustomCardHandler(w http.ResponseWriter, r *http.Request) {
	weekIDParam := chi.URLParam(r, "week_id")
	weekID, ok := parseUUID(w, weekIDParam)
	if !ok {
		return
	}

	userIDStr := getUserID(r)
	userID, ok := parseUUID(w, userIDStr)
	if !ok {
		return
	}

	var req content.CreateCustomCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", "error", err)
		ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	card, err := s.contentSrv.CreateCustomCard(r.Context(), userID, weekID, req.Front, req.Back)
	if err != nil {
		if strings.Contains(err.Error(), "empty") {
			ResponseWithErr(w, http.StatusBadRequest, err.Error())
			return
		}
		slog.Error("failed to create custom card", "error", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to create custom card")
		return
	}

	ResponseWithJSON(w, http.StatusCreated, card)
}

// GetUserDeckHandler retrieves all cards in user's deck for a specific week
func (s *HTTPServer) GetUserDeckHandler(w http.ResponseWriter, r *http.Request) {
	weekIDParam := chi.URLParam(r, "week_id")
	weekID, ok := parseUUID(w, weekIDParam)
	if !ok {
		return
	}

	userIDStr := getUserID(r)
	userID, ok := parseUUID(w, userIDStr)
	if !ok {
		return
	}

	cards, err := s.contentSrv.GetUserDeckForWeek(r.Context(), userID, weekID)
	if err != nil {
		slog.Error("failed to get user deck", "error", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to get deck")
		return
	}

	ResponseWithJSON(w, http.StatusOK, cards)
}

// UpdateDeckCardHandler updates a card in user's deck
func (s *HTTPServer) UpdateDeckCardHandler(w http.ResponseWriter, r *http.Request) {
	cardIDParam := chi.URLParam(r, "card_id")
	cardID, ok := parseUUID(w, cardIDParam)
	if !ok {
		return
	}

	userIDStr := getUserID(r)
	userID, ok := parseUUID(w, userIDStr)
	if !ok {
		return
	}

	var req content.UpdateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", "error", err)
		ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := s.contentSrv.UpdateDeckCard(r.Context(), cardID, userID, req.Front, req.Back)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "access denied") {
			ResponseWithErr(w, http.StatusNotFound, err.Error())
			return
		}
		slog.Error("failed to update card", "error", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to update card")
		return
	}

	ResponseWithJSON(w, http.StatusOK, nil)
}

// RemoveDeckCardHandler removes a card from user's deck
func (s *HTTPServer) RemoveDeckCardHandler(w http.ResponseWriter, r *http.Request) {
	cardIDParam := chi.URLParam(r, "card_id")
	cardID, ok := parseUUID(w, cardIDParam)
	if !ok {
		return
	}

	userIDStr := getUserID(r)
	userID, ok := parseUUID(w, userIDStr)
	if !ok {
		return
	}

	err := s.contentSrv.RemoveCardFromDeck(r.Context(), cardID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "access denied") {
			ResponseWithErr(w, http.StatusNotFound, err.Error())
			return
		}
		slog.Error("failed to remove card", "error", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to remove card")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RecordCardReviewHandler records a review session for a card
func (s *HTTPServer) RecordCardReviewHandler(w http.ResponseWriter, r *http.Request) {
	cardIDParam := chi.URLParam(r, "card_id")
	cardID, ok := parseUUID(w, cardIDParam)
	if !ok {
		return
	}

	userIDStr := getUserID(r)
	userID, ok := parseUUID(w, userIDStr)
	if !ok {
		return
	}

	var req content.RecordReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", "error", err)
		ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := s.contentSrv.RecordCardReview(r.Context(), cardID, userID, req.DifficultyRating)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "access denied") {
			ResponseWithErr(w, http.StatusNotFound, err.Error())
			return
		}
		if strings.Contains(err.Error(), "between 1 and 5") {
			ResponseWithErr(w, http.StatusBadRequest, err.Error())
			return
		}
		slog.Error("failed to record review", "error", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to record review")
		return
	}

	ResponseWithJSON(w, http.StatusOK, nil)
}

// GetDeckStatsHandler retrieves statistics for user's deck in a week
func (s *HTTPServer) GetDeckStatsHandler(w http.ResponseWriter, r *http.Request) {
	weekIDParam := chi.URLParam(r, "week_id")
	weekID, ok := parseUUID(w, weekIDParam)
	if !ok {
		return
	}

	userIDStr := getUserID(r)
	userID, ok := parseUUID(w, userIDStr)
	if !ok {
		return
	}

	stats, err := s.contentSrv.GetDeckStats(r.Context(), userID, weekID)
	if err != nil {
		slog.Error("failed to get deck stats", "error", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to get statistics")
		return
	}

	ResponseWithJSON(w, http.StatusOK, stats)
}
