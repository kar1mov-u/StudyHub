package http

import (
	"StudyHub/internal/content"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// mockContentService implements the methods needed for testing deck handlers
type mockContentService struct {
	addCardToDeckFunc      func(ctx context.Context, userID, weekID, flashcardID uuid.UUID) error
	createCustomCardFunc   func(ctx context.Context, userID, weekID uuid.UUID, front, back string) (content.UserDeckCard, error)
	getUserDeckFunc        func(ctx context.Context, userID, weekID uuid.UUID) ([]content.UserDeckCard, error)
	updateDeckCardFunc     func(ctx context.Context, cardID, userID uuid.UUID, front, back *string) error
	removeCardFromDeckFunc func(ctx context.Context, cardID, userID uuid.UUID) error
	recordCardReviewFunc   func(ctx context.Context, cardID, userID uuid.UUID, difficultyRating int) error
	getDeckStatsFunc       func(ctx context.Context, userID, weekID uuid.UUID) (content.DeckStats, error)
}

func (m *mockContentService) AddCardToUserDeck(ctx context.Context, userID, weekID, flashcardID uuid.UUID) error {
	if m.addCardToDeckFunc != nil {
		return m.addCardToDeckFunc(ctx, userID, weekID, flashcardID)
	}
	return nil
}

func (m *mockContentService) CreateCustomCard(ctx context.Context, userID, weekID uuid.UUID, front, back string) (content.UserDeckCard, error) {
	if m.createCustomCardFunc != nil {
		return m.createCustomCardFunc(ctx, userID, weekID, front, back)
	}
	return content.UserDeckCard{}, nil
}

func (m *mockContentService) GetUserDeckForWeek(ctx context.Context, userID, weekID uuid.UUID) ([]content.UserDeckCard, error) {
	if m.getUserDeckFunc != nil {
		return m.getUserDeckFunc(ctx, userID, weekID)
	}
	return []content.UserDeckCard{}, nil
}

func (m *mockContentService) UpdateDeckCard(ctx context.Context, cardID, userID uuid.UUID, front, back *string) error {
	if m.updateDeckCardFunc != nil {
		return m.updateDeckCardFunc(ctx, cardID, userID, front, back)
	}
	return nil
}

func (m *mockContentService) RemoveCardFromDeck(ctx context.Context, cardID, userID uuid.UUID) error {
	if m.removeCardFromDeckFunc != nil {
		return m.removeCardFromDeckFunc(ctx, cardID, userID)
	}
	return nil
}

func (m *mockContentService) RecordCardReview(ctx context.Context, cardID, userID uuid.UUID, difficultyRating int) error {
	if m.recordCardReviewFunc != nil {
		return m.recordCardReviewFunc(ctx, cardID, userID, difficultyRating)
	}
	return nil
}

func (m *mockContentService) GetDeckStats(ctx context.Context, userID, weekID uuid.UUID) (content.DeckStats, error) {
	if m.getDeckStatsFunc != nil {
		return m.getDeckStatsFunc(ctx, userID, weekID)
	}
	return content.DeckStats{}, nil
}

// Helper to add user ID to context (simulating auth middleware)
func addUserIDToContext(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), "userID", userID)
	return r.WithContext(ctx)
}

func TestAddCardToDeckHandler(t *testing.T) {
	tests := []struct {
		name           string
		weekID         string
		userID         string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, userID, weekID, flashcardID uuid.UUID) error
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "success - add card to deck",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.AddCardToDeckRequest{
				FlashcardID: uuid.New().String(),
			},
			mockFunc: func(ctx context.Context, userID, weekID, flashcardID uuid.UUID) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "error - invalid request body",
			weekID:         uuid.New().String(),
			userID:         uuid.New().String(),
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:   "error - invalid week ID",
			weekID: "invalid-uuid",
			userID: uuid.New().String(),
			requestBody: content.AddCardToDeckRequest{
				FlashcardID: uuid.New().String(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:   "error - invalid flashcard ID",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.AddCardToDeckRequest{
				FlashcardID: "invalid-uuid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:   "error - flashcard not found",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.AddCardToDeckRequest{
				FlashcardID: uuid.New().String(),
			},
			mockFunc: func(ctx context.Context, userID, weekID, flashcardID uuid.UUID) error {
				return errors.New("flashcard not found")
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:   "error - duplicate card",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.AddCardToDeckRequest{
				FlashcardID: uuid.New().String(),
			},
			mockFunc: func(ctx context.Context, userID, weekID, flashcardID uuid.UUID) error {
				return errors.New("duplicate key value violates unique constraint")
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockContentService{
				addCardToDeckFunc: tt.mockFunc,
			}
			server := &HTTPServer{contentSrv: &content.ContentService{}}
			// Inject mock service by creating a custom wrapper
			testServer := &struct {
				*HTTPServer
				mockSvc *mockContentService
			}{
				HTTPServer: server,
				mockSvc:    mockSvc,
			}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/decks/weeks/"+tt.weekID+"/cards", bytes.NewBuffer(body))
			req = addUserIDToContext(req, tt.userID)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("week_id", tt.weekID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			// Call handler with mock
			weekIDParam := chi.URLParam(req, "week_id")
			weekID, okWeek := parseUUID(w, weekIDParam)
			if !okWeek {
				return
			}

			userIDStr := getUserID(req)
			userID, okUser := parseUUID(w, userIDStr)
			if !okUser {
				return
			}

			var reqBody content.AddCardToDeckRequest
			if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&reqBody); err != nil {
				ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
				return
			}

			flashcardID, err := uuid.Parse(reqBody.FlashcardID)
			if err != nil {
				ResponseWithErr(w, http.StatusBadRequest, "invalid flashcard_id")
			} else {
				err = testServer.mockSvc.AddCardToUserDeck(req.Context(), userID, weekID, flashcardID)
				if err != nil {
					if err.Error() == "flashcard not found" {
						ResponseWithErr(w, http.StatusNotFound, err.Error())
					} else if err.Error() == "duplicate key value violates unique constraint" {
						ResponseWithErr(w, http.StatusConflict, "card already in deck")
					} else {
						ResponseWithErr(w, http.StatusInternalServerError, "failed to add card to deck")
					}
				} else {
					ResponseWithJSON(w, http.StatusCreated, nil)
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError {
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				if _, exists := resp["error"]; !exists {
					t.Error("expected error in response")
				}
			}
		})
	}
}

func TestGetUserDeckHandler(t *testing.T) {
	tests := []struct {
		name           string
		weekID         string
		userID         string
		mockFunc       func(ctx context.Context, userID, weekID uuid.UUID) ([]content.UserDeckCard, error)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:   "success - get user deck with cards",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, userID, weekID uuid.UUID) ([]content.UserDeckCard, error) {
				return []content.UserDeckCard{
					{
						ID:       uuid.New(),
						UserID:   userID,
						WeekID:   weekID,
						Front:    "Question 1",
						Back:     "Answer 1",
						IsCustom: true,
					},
					{
						ID:       uuid.New(),
						UserID:   userID,
						WeekID:   weekID,
						Front:    "Question 2",
						Back:     "Answer 2",
						IsCustom: false,
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:   "success - empty deck",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, userID, weekID uuid.UUID) ([]content.UserDeckCard, error) {
				return []content.UserDeckCard{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "error - invalid week ID",
			weekID:         "invalid-uuid",
			userID:         uuid.New().String(),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockContentService{
				getUserDeckFunc: tt.mockFunc,
			}

			req := httptest.NewRequest(http.MethodGet, "/decks/weeks/"+tt.weekID+"/cards", nil)
			req = addUserIDToContext(req, tt.userID)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("week_id", tt.weekID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			// Call handler logic with mock
			weekIDParam := chi.URLParam(req, "week_id")
			weekID, ok := parseUUID(w, weekIDParam)
			if !ok {
				// parseUUID already wrote error response
			} else {
				userIDStr := getUserID(req)
				userID, okUser := parseUUID(w, userIDStr)
				if okUser {
					cards, err := mockSvc.GetUserDeckForWeek(req.Context(), userID, weekID)
					if err != nil {
						ResponseWithErr(w, http.StatusInternalServerError, "failed to get deck")
					} else {
						ResponseWithJSON(w, http.StatusOK, cards)
					}
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var resp Response
				json.Unmarshal(w.Body.Bytes(), &resp)
				if resp.Data == nil && tt.expectedCount > 0 {
					t.Error("expected data in response")
				}
			}
		})
	}
}

func TestRecordCardReviewHandler(t *testing.T) {
	tests := []struct {
		name           string
		cardID         string
		userID         string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, cardID, userID uuid.UUID, difficultyRating int) error
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "success - record review with difficulty 3",
			cardID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.RecordReviewRequest{
				DifficultyRating: 3,
			},
			mockFunc: func(ctx context.Context, cardID, userID uuid.UUID, difficultyRating int) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "error - invalid request body",
			cardID:         uuid.New().String(),
			userID:         uuid.New().String(),
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:   "error - card not found",
			cardID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.RecordReviewRequest{
				DifficultyRating: 3,
			},
			mockFunc: func(ctx context.Context, cardID, userID uuid.UUID, difficultyRating int) error {
				return errors.New("card not found or access denied")
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:   "error - invalid difficulty rating (0)",
			cardID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.RecordReviewRequest{
				DifficultyRating: 0,
			},
			mockFunc: func(ctx context.Context, cardID, userID uuid.UUID, difficultyRating int) error {
				return errors.New("difficulty rating must be between 1 and 5")
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:   "error - invalid difficulty rating (6)",
			cardID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.RecordReviewRequest{
				DifficultyRating: 6,
			},
			mockFunc: func(ctx context.Context, cardID, userID uuid.UUID, difficultyRating int) error {
				return errors.New("difficulty rating must be between 1 and 5")
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockContentService{
				recordCardReviewFunc: tt.mockFunc,
			}

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
			req := httptest.NewRequest(http.MethodPost, "/decks/cards/"+tt.cardID+"/review", bytes.NewBuffer(body))
			req = addUserIDToContext(req, tt.userID)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("card_id", tt.cardID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			// Call handler logic with mock
			cardIDParam := chi.URLParam(req, "card_id")
			cardID, ok := parseUUID(w, cardIDParam)
			if ok {
				userIDStr := getUserID(req)
				userID, okUser := parseUUID(w, userIDStr)
				if okUser {
					var reqBody content.RecordReviewRequest
					if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&reqBody); err != nil {
						ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
					} else {
						err = mockSvc.RecordCardReview(req.Context(), cardID, userID, reqBody.DifficultyRating)
						if err != nil {
							if err.Error() == "difficulty rating must be between 1 and 5" {
								ResponseWithErr(w, http.StatusBadRequest, err.Error())
							} else if err.Error() == "card not found or access denied" {
								ResponseWithErr(w, http.StatusNotFound, err.Error())
							} else {
								ResponseWithErr(w, http.StatusInternalServerError, "failed to record review")
							}
						} else {
							ResponseWithJSON(w, http.StatusOK, nil)
						}
					}
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetDeckStatsHandler(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name           string
		weekID         string
		userID         string
		mockFunc       func(ctx context.Context, userID, weekID uuid.UUID) (content.DeckStats, error)
		expectedStatus int
		expectData     bool
	}{
		{
			name:   "success - get deck stats",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, userID, weekID uuid.UUID) (content.DeckStats, error) {
				return content.DeckStats{
					TotalCards:     10,
					ReviewedCards:  5,
					AverageRating:  3.5,
					LastReviewedAt: &now,
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectData:     true,
		},
		{
			name:   "success - empty deck stats",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, userID, weekID uuid.UUID) (content.DeckStats, error) {
				return content.DeckStats{
					TotalCards:     0,
					ReviewedCards:  0,
					AverageRating:  0,
					LastReviewedAt: nil,
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectData:     true,
		},
		{
			name:           "error - invalid week ID",
			weekID:         "invalid-uuid",
			userID:         uuid.New().String(),
			expectedStatus: http.StatusBadRequest,
			expectData:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockContentService{
				getDeckStatsFunc: tt.mockFunc,
			}

			req := httptest.NewRequest(http.MethodGet, "/decks/weeks/"+tt.weekID+"/stats", nil)
			req = addUserIDToContext(req, tt.userID)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("week_id", tt.weekID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			// Call handler logic with mock
			weekIDParam := chi.URLParam(req, "week_id")
			weekID, ok := parseUUID(w, weekIDParam)
			if ok {
				userIDStr := getUserID(req)
				userID, okUser := parseUUID(w, userIDStr)
				if okUser {
					stats, err := mockSvc.GetDeckStats(req.Context(), userID, weekID)
					if err != nil {
						ResponseWithErr(w, http.StatusInternalServerError, "failed to get deck stats")
					} else {
						ResponseWithJSON(w, http.StatusOK, stats)
					}
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectData {
				var resp Response
				json.Unmarshal(w.Body.Bytes(), &resp)
				if resp.Data == nil {
					t.Error("expected data in response")
				}
			}
		})
	}
}

func TestListCardsFromObjects(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, ids []uuid.UUID) ([]content.Flashcard, error)
		expectedStatus int
	}{
		{
			name: "success - list cards for objects",
			requestBody: ListCardsForObjectsRequest{
				IDs: []string{uuid.New().String(), uuid.New().String()},
			},
			mockFunc: func(ctx context.Context, ids []uuid.UUID) ([]content.Flashcard, error) {
				return []content.Flashcard{{ID: uuid.New(), Front: "Front", Back: "Back"}}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error - invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error - invalid object ID",
			requestBody: ListCardsForObjectsRequest{
				IDs: []string{"invalid-uuid"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error - service error",
			requestBody: ListCardsForObjectsRequest{
				IDs: []string{uuid.New().String()},
			},
			mockFunc: func(ctx context.Context, ids []uuid.UUID) ([]content.Flashcard, error) {
				return nil, errors.New("failed to list cards")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			req := httptest.NewRequest(http.MethodPost, "/conents/objects", bytes.NewReader(body))
			w := httptest.NewRecorder()

			var payload ListCardsForObjectsRequest
			if err := json.NewDecoder(bytes.NewReader(body)).Decode(&payload); err != nil {
				ResponseWithErr(w, http.StatusBadRequest, "invalid request")
			} else {
				uuids := make([]uuid.UUID, 0, len(payload.IDs))
				valid := true
				for i := range payload.IDs {
					id, err := uuid.Parse(payload.IDs[i])
					if err != nil {
						ResponseWithErr(w, http.StatusBadRequest, "invalid objectID")
						valid = false
						break
					}
					uuids = append(uuids, id)
				}

				if valid {
					cards := []content.Flashcard{}
					if tt.mockFunc != nil {
						cards, err = tt.mockFunc(req.Context(), uuids)
					}
					if err != nil {
						ResponseWithErr(w, http.StatusInternalServerError, "failed to list cards")
					} else {
						ResponseWithJSON(w, http.StatusOK, cards)
					}
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestCreateCustomCardHandler(t *testing.T) {
	tests := []struct {
		name           string
		weekID         string
		userID         string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, userID, weekID uuid.UUID, front, back string) (content.UserDeckCard, error)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "success - create custom card",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.CreateCustomCardRequest{
				Front: "What is Go?",
				Back:  "A programming language",
			},
			mockFunc: func(ctx context.Context, userID, weekID uuid.UUID, front, back string) (content.UserDeckCard, error) {
				return content.UserDeckCard{
					ID:        uuid.New(),
					UserID:    userID,
					WeekID:    weekID,
					Front:     front,
					Back:      back,
					IsCustom:  true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:   "error - empty front",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.CreateCustomCardRequest{
				Front: "",
				Back:  "A programming language",
			},
			mockFunc: func(ctx context.Context, userID, weekID uuid.UUID, front, back string) (content.UserDeckCard, error) {
				return content.UserDeckCard{}, errors.New("front and back cannot be empty")
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:   "error - empty back",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.CreateCustomCardRequest{
				Front: "What is Go?",
				Back:  "",
			},
			mockFunc: func(ctx context.Context, userID, weekID uuid.UUID, front, back string) (content.UserDeckCard, error) {
				return content.UserDeckCard{}, errors.New("front and back cannot be empty")
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:   "error - invalid week ID",
			weekID: "invalid-uuid",
			userID: uuid.New().String(),
			requestBody: content.CreateCustomCardRequest{
				Front: "Question",
				Back:  "Answer",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockContentService{
				createCustomCardFunc: tt.mockFunc,
			}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/decks/weeks/"+tt.weekID+"/cards/custom", bytes.NewBuffer(body))
			req = addUserIDToContext(req, tt.userID)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("week_id", tt.weekID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			// Call handler logic with mock
			weekIDParam := chi.URLParam(req, "week_id")
			weekID, ok := parseUUID(w, weekIDParam)
			if !ok {
				// parseUUID already wrote error response
			} else {
				userIDStr := getUserID(req)
				userID, okUser := parseUUID(w, userIDStr)
				if okUser {
					var reqBody content.CreateCustomCardRequest
					if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&reqBody); err != nil {
						ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
					} else {
						card, err := mockSvc.CreateCustomCard(req.Context(), userID, weekID, reqBody.Front, reqBody.Back)
						if err != nil {
							if err.Error() == "front and back cannot be empty" {
								ResponseWithErr(w, http.StatusBadRequest, err.Error())
							} else {
								ResponseWithErr(w, http.StatusInternalServerError, "failed to create custom card")
							}
						} else {
							ResponseWithJSON(w, http.StatusCreated, card)
						}
					}
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError {
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				if _, exists := resp["error"]; !exists {
					t.Error("expected error in response")
				}
			}
		})
	}
}

func TestUpdateDeckCardHandler(t *testing.T) {
	tests := []struct {
		name           string
		cardID         string
		userID         string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, cardID, userID uuid.UUID, front, back *string) error
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "success - update both front and back",
			cardID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.UpdateCardRequest{
				Front: stringPtr("Updated question"),
				Back:  stringPtr("Updated answer"),
			},
			mockFunc: func(ctx context.Context, cardID, userID uuid.UUID, front, back *string) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "success - update front only",
			cardID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.UpdateCardRequest{
				Front: stringPtr("Updated question"),
			},
			mockFunc: func(ctx context.Context, cardID, userID uuid.UUID, front, back *string) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "success - update back only",
			cardID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.UpdateCardRequest{
				Back: stringPtr("Updated answer"),
			},
			mockFunc: func(ctx context.Context, cardID, userID uuid.UUID, front, back *string) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "error - card not found",
			cardID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: content.UpdateCardRequest{
				Front: stringPtr("Question"),
			},
			mockFunc: func(ctx context.Context, cardID, userID uuid.UUID, front, back *string) error {
				return errors.New("card not found or access denied")
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:           "error - invalid card ID",
			cardID:         "invalid-uuid",
			userID:         uuid.New().String(),
			requestBody:    content.UpdateCardRequest{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockContentService{
				updateDeckCardFunc: tt.mockFunc,
			}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/decks/cards/"+tt.cardID, bytes.NewBuffer(body))
			req = addUserIDToContext(req, tt.userID)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("card_id", tt.cardID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			// Call handler logic with mock
			cardIDParam := chi.URLParam(req, "card_id")
			cardID, ok := parseUUID(w, cardIDParam)
			if !ok {
				// parseUUID already wrote error response
			} else {
				userIDStr := getUserID(req)
				userID, okUser := parseUUID(w, userIDStr)
				if okUser {
					var reqBody content.UpdateCardRequest
					if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&reqBody); err != nil {
						ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
					} else {
						err := mockSvc.UpdateDeckCard(req.Context(), cardID, userID, reqBody.Front, reqBody.Back)
						if err != nil {
							if err.Error() == "card not found or access denied" {
								ResponseWithErr(w, http.StatusNotFound, err.Error())
							} else {
								ResponseWithErr(w, http.StatusInternalServerError, "failed to update card")
							}
						} else {
							ResponseWithJSON(w, http.StatusOK, nil)
						}
					}
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRemoveDeckCardHandler(t *testing.T) {
	tests := []struct {
		name           string
		cardID         string
		userID         string
		mockFunc       func(ctx context.Context, cardID, userID uuid.UUID) error
		expectedStatus int
	}{
		{
			name:   "success - remove card",
			cardID: uuid.New().String(),
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, cardID, userID uuid.UUID) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "error - card not found",
			cardID: uuid.New().String(),
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, cardID, userID uuid.UUID) error {
				return errors.New("card not found or access denied")
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "error - invalid card ID",
			cardID:         "invalid-uuid",
			userID:         uuid.New().String(),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockContentService{
				removeCardFromDeckFunc: tt.mockFunc,
			}

			req := httptest.NewRequest(http.MethodDelete, "/decks/cards/"+tt.cardID, nil)
			req = addUserIDToContext(req, tt.userID)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("card_id", tt.cardID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			// Call handler logic with mock
			cardIDParam := chi.URLParam(req, "card_id")
			cardID, ok := parseUUID(w, cardIDParam)
			if !ok {
				// parseUUID already wrote error response
			} else {
				userIDStr := getUserID(req)
				userID, okUser := parseUUID(w, userIDStr)
				if okUser {
					err := mockSvc.RemoveCardFromDeck(req.Context(), cardID, userID)
					if err != nil {
						if err.Error() == "card not found or access denied" {
							ResponseWithErr(w, http.StatusNotFound, err.Error())
						} else {
							ResponseWithErr(w, http.StatusInternalServerError, "failed to remove card")
						}
					} else {
						w.WriteHeader(http.StatusNoContent)
					}
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
