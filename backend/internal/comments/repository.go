package comments

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRespositoryPostres struct {
	pool *pgxpool.Pool
}

func NewCommentRepositoryPostgres(pool *pgxpool.Pool) *CommentRespositoryPostres {
	return &CommentRespositoryPostres{
		pool: pool,
	}
}

func (r *CommentRespositoryPostres) CreateComment(comment Comment) error {
	_, err := r.pool.Exec(context.Background(), "INSERT INTO week_comments (id, week_id, user_id, reply, content) VALUES ($1, $2, $3, $4, $5)",
		comment.ID, comment.WeekID, comment.UserID, comment.ReplyID, comment.Content)
	return err
}

func (r *CommentRespositoryPostres) GetCommentsByWeekID(weekID string) ([]Comment, error) {
	rows, err := r.pool.Query(context.Background(), "SELECT id, week_id, user_id, reply, content, upvote, downvote, created_at FROM week_comments WHERE week_id = $1", weekID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.WeekID, &comment.UserID, &comment.ReplyID, &comment.Content, &comment.Upvote, &comment.Downvote, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRespositoryPostres) UpvoteComment(commentID string) error {
	_, err := r.pool.Exec(context.Background(), "UPDATE week_comments SET upvote = upvote + 1 WHERE id = $1", commentID)
	return err
}

func (r *CommentRespositoryPostres) DownvoteComment(commentID string) error {
	_, err := r.pool.Exec(context.Background(), "UPDATE week_comments SET downvote = downvote + 1 WHERE id = $1", commentID)
	return err
}
