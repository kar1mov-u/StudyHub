package comments

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	WeekID    uuid.UUID  `json:"week_id"`
	ReplyID   *uuid.UUID `json:"reply_id,omitempty"`
	Content   string     `json:"content"`
	Upvote    int        `json:"upvote"`
	Downvote  int        `json:"downvote"`
	CreatedAt time.Time  `json:"created_at"`
}
