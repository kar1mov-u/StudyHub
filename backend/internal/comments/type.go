package comments

import "github.com/google/uuid"

type Comment struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	WeekID    uuid.UUID  `json:"week_id"`
	ReplyID   *uuid.UUID `json:"reply_id,omitempty"`
	Content   string     `json:"content"`
	Upvote    int        `json:"upvote"`
	Downvote  int        `json:"downvote"`
	CreatedAt int64      `json:"created_at"`
}
