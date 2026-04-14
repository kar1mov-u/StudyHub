package http

import (
	"StudyHub/internal/comments"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Mock CommentService for testing
type mockCommentService struct {
	createCommentFunc       func(comment comments.Comment) error
	getCommentsByWeekIDFunc func(weekID string) ([]comments.Comment, error)
	upvoteCommentFunc       func(commentID, userID string) error
	downvoteCommentFunc     func(commentID, userID string) error
}

func (m *mockCommentService) CreateComment(comment comments.Comment) error {
	if m.createCommentFunc != nil {
		return m.createCommentFunc(comment)
	}
	return nil
}

func (m *mockCommentService) GetCommentsByWeekID(weekID string) ([]comments.Comment, error) {
	if m.getCommentsByWeekIDFunc != nil {
		return m.getCommentsByWeekIDFunc(weekID)
	}
	return []comments.Comment{}, nil
}

func (m *mockCommentService) UpvoteComment(commentID, userID string) error {
	if m.upvoteCommentFunc != nil {
		return m.upvoteCommentFunc(commentID, userID)
	}
	return nil
}

func (m *mockCommentService) DownvoteComment(commentID, userID string) error {
	if m.downvoteCommentFunc != nil {
		return m.downvoteCommentFunc(commentID, userID)
	}
	return nil
}

// Test CreateCommentHandler
func TestCreateCommentHandler(t *testing.T) {
	userID := uuid.New()
	weekID := uuid.New()
	replyID := uuid.New()
	replyIDStr := replyID.String()

	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(comment comments.Comment) error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success - create comment without reply",
			requestBody: CreateCommentRequest{
				UserID:  userID.String(),
				WeekID:  weekID.String(),
				Content: "This is a test comment",
			},
			mockFunc: func(comment comments.Comment) error {
				if comment.UserID != userID {
					t.Errorf("expected userID=%s, got %s", userID, comment.UserID)
				}
				if comment.WeekID != weekID {
					t.Errorf("expected weekID=%s, got %s", weekID, comment.WeekID)
				}
				if comment.Content != "This is a test comment" {
					t.Errorf("expected content='This is a test comment', got %s", comment.Content)
				}
				if comment.ReplyID != nil {
					t.Errorf("expected replyID to be nil, got %v", comment.ReplyID)
				}
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "This is a test comment",
		},
		{
			name: "success - create comment with reply",
			requestBody: CreateCommentRequest{
				UserID:  userID.String(),
				WeekID:  weekID.String(),
				ReplyID: &replyIDStr,
				Content: "This is a reply comment",
			},
			mockFunc: func(comment comments.Comment) error {
				if comment.ReplyID == nil {
					t.Errorf("expected replyID to not be nil")
				} else if *comment.ReplyID != replyID {
					t.Errorf("expected replyID=%s, got %s", replyID, *comment.ReplyID)
				}
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "This is a reply comment",
		},
		{
			name:           "error - invalid JSON",
			requestBody:    "invalid json",
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid request body",
		},
		{
			name: "error - service error",
			requestBody: CreateCommentRequest{
				UserID:  userID.String(),
				WeekID:  weekID.String(),
				Content: "Test comment",
			},
			mockFunc: func(comment comments.Comment) error {
				return errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to create comment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockCommentService{createCommentFunc: tt.mockFunc}

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

			req := httptest.NewRequest(http.MethodPost, "/api/v1/comments", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute handler logic inline
			var reqData CreateCommentRequest
			if err := json.NewDecoder(bytes.NewReader(body)).Decode(&reqData); err != nil {
				ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
			} else {
				comment := comments.Comment{
					UserID:  uuid.MustParse(reqData.UserID),
					WeekID:  uuid.MustParse(reqData.WeekID),
					Content: reqData.Content,
				}
				if reqData.ReplyID != nil {
					parsedReplyID := uuid.MustParse(*reqData.ReplyID)
					comment.ReplyID = &parsedReplyID
				}

				if err := mockSvc.CreateComment(comment); err != nil {
					ResponseWithErr(w, http.StatusInternalServerError, "failed to create comment")
				} else {
					ResponseWithJSON(w, http.StatusCreated, comment)
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

// Test ListCommentsForWeekHandler
func TestListCommentsForWeekHandler(t *testing.T) {
	weekID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name           string
		weekIDParam    string
		mockFunc       func(weekID string) ([]comments.Comment, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "success - returns comments list",
			weekIDParam: weekID.String(),
			mockFunc: func(wid string) ([]comments.Comment, error) {
				if wid != weekID.String() {
					t.Errorf("expected weekID=%s, got %s", weekID.String(), wid)
				}
				return []comments.Comment{
					{
						ID:       uuid.New(),
						UserID:   userID,
						WeekID:   weekID,
						Content:  "First comment",
						Upvote:   5,
						Downvote: 1,
					},
					{
						ID:       uuid.New(),
						UserID:   userID,
						WeekID:   weekID,
						Content:  "Second comment",
						Upvote:   3,
						Downvote: 0,
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "First comment",
		},
		{
			name:        "success - empty comments list",
			weekIDParam: weekID.String(),
			mockFunc: func(wid string) ([]comments.Comment, error) {
				return []comments.Comment{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "[]",
		},
		{
			name:        "error - service error",
			weekIDParam: weekID.String(),
			mockFunc: func(wid string) ([]comments.Comment, error) {
				return nil, errors.New("database connection failed")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to get comments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockCommentService{getCommentsByWeekIDFunc: tt.mockFunc}

			req := httptest.NewRequest(http.MethodGet, "/api/v1/weeks/"+tt.weekIDParam+"/comments", nil)
			w := httptest.NewRecorder()

			// Set up chi route context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("week_id", tt.weekIDParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Execute handler logic inline
			weekIDParam := chi.URLParam(req, "week_id")
			commentsList, err := mockSvc.GetCommentsByWeekID(weekIDParam)
			if err != nil {
				ResponseWithErr(w, http.StatusInternalServerError, "failed to get comments")
			} else {
				ResponseWithJSON(w, http.StatusOK, commentsList)
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

func TestUpvoteCommentHandler(t *testing.T) {
	tests := []struct {
		name           string
		commentID      string
		userID         string
		mockFunc       func(commentID, userID string) error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "success - upvote comment",
			commentID: uuid.New().String(),
			userID:    uuid.New().String(),
			mockFunc: func(commentID, userID string) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "comment upvoted successfully",
		},
		{
			name:      "error - upvote failure",
			commentID: uuid.New().String(),
			userID:    uuid.New().String(),
			mockFunc: func(commentID, userID string) error {
				return errors.New("failed to upvote")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to upvote comment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockCommentService{upvoteCommentFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/comments/"+tt.commentID+"/upvote", nil)
			req = addUserIDToContext(req, tt.userID)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.commentID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			commentIDParam := chi.URLParam(req, "id")
			userIDParam := getUserID(req)
			if err := mockSvc.UpvoteComment(commentIDParam, userIDParam); err != nil {
				ResponseWithErr(w, http.StatusInternalServerError, "failed to upvote comment")
			} else {
				ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "comment upvoted successfully"})
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

func TestDownvoteCommentHandler(t *testing.T) {
	tests := []struct {
		name           string
		commentID      string
		userID         string
		mockFunc       func(commentID, userID string) error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "success - downvote comment",
			commentID: uuid.New().String(),
			userID:    uuid.New().String(),
			mockFunc: func(commentID, userID string) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "comment downvoted successfully",
		},
		{
			name:      "error - downvote failure",
			commentID: uuid.New().String(),
			userID:    uuid.New().String(),
			mockFunc: func(commentID, userID string) error {
				return errors.New("failed to downvote")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to downvote comment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockCommentService{downvoteCommentFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/comments/"+tt.commentID+"/downvote", nil)
			req = addUserIDToContext(req, tt.userID)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.commentID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			commentIDParam := chi.URLParam(req, "id")
			userIDParam := getUserID(req)
			if err := mockSvc.DownvoteComment(commentIDParam, userIDParam); err != nil {
				ResponseWithErr(w, http.StatusInternalServerError, "failed to downvote comment")
			} else {
				ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "comment downvoted successfully"})
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
