package http

import (
	"StudyHub/internal/comments"
	"encoding/json"
	"net/http"

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
