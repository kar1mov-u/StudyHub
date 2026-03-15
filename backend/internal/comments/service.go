package comments

import "github.com/google/uuid"

type CommentRepository interface {
	CreateComment(comment Comment) error
	GetCommentsByWeekID(weekID string) ([]Comment, error)
	UpvoteComment(commentID string) error
	DownvoteComment(commentID string) error
}

type CommentService struct {
	repo CommentRepository
}

func NewCommentService(repo CommentRepository) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

func (s *CommentService) CreateComment(comment Comment) error {
	comment.ID = uuid.New()
	return s.repo.CreateComment(comment)
}

func (s *CommentService) GetCommentsByWeekID(weekID string) ([]Comment, error) {
	return s.repo.GetCommentsByWeekID(weekID)
}

func (s *CommentService) UpvoteComment(commentID string) error {
	return s.repo.UpvoteComment(commentID)
}

func (s *CommentService) DownvoteComment(commentID string) error {
	return s.repo.DownvoteComment(commentID)
}
