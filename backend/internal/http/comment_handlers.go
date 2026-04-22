package http

import (
	"StudyHub/internal/comments"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateCommentRequest struct {
	UserID  string  `json:"user_id"`
	WeekID  string  `json:"week_id"`
	ReplyID *string `json:"reply_id,omitempty"`
	Content string  `json:"content"`
}

func (s *HTTPServer) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	//parse the request body to get the comment data
	var req CreateCommentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
		return
	}
	comment := comments.Comment{
		UserID:  uuid.MustParse(req.UserID),
		WeekID:  uuid.MustParse(req.WeekID),
		Content: req.Content,
	}
	if req.ReplyID != nil {
		replyID := uuid.MustParse(*req.ReplyID)
		comment.ReplyID = &replyID
	}

	err = s.commentSrv.CreateComment(comment)
	if err != nil {
		ResponseWithErr(w, http.StatusInternalServerError, "failed to create comment")
		return
	}
	ResponseWithJSON(w, http.StatusCreated, comment)
}

func (s *HTTPServer) ListCommentsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	weekID := chi.URLParam(r, "week_id")
	comments, err := s.commentSrv.GetCommentsByWeekID(weekID)
	if err != nil {
		slog.Error("failed to get comments for week")
		ResponseWithErr(w, http.StatusInternalServerError, "failed to get comments")
		return
	}
	ResponseWithJSON(w, http.StatusOK, comments)
}

func (s *HTTPServer) UpvoteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "id")
	userID := getUserID(r)
	err := s.commentSrv.UpvoteComment(commentID, userID)
	if err != nil {
		slog.Error("failed to upvote comment")
		ResponseWithErr(w, http.StatusInternalServerError, "failed to upvote comment")
		return
	}
	ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "comment upvoted successfully"})
}

func (s *HTTPServer) DownvoteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "id")
	userID := getUserID(r)
	err := s.commentSrv.DownvoteComment(commentID, userID)
	if err != nil {
		slog.Error("failed to downvote comment")
		ResponseWithErr(w, http.StatusInternalServerError, "failed to downvote comment")
		return
	}
	ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "comment downvoted successfully"})
}
