package comments

import "github.com/google/uuid"

type CommentRepository interface {
	CreateComment(comment Comment) error
	GetCommentsByWeekID(weekID string) ([]Comment, error)
	UpvoteComment(commentID string) error
	DownvoteComment(commentID string) error
	UserHasVoted(commentID, userID string) (bool, error)
	CreateVote(commentID, userID string, isUpvote bool) error
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

func (s *CommentService) UpvoteComment(commentID, userID string) error {
	hasVoted, err := s.repo.UserHasVoted(commentID, userID)
	if err != nil {
		return err
	}
	if hasVoted {
		return nil
	}
	err = s.repo.CreateVote(commentID, userID, true)
	if err != nil {
		return err
	}
	return s.repo.UpvoteComment(commentID)
}

func (s *CommentService) DownvoteComment(commentID, userID string) error {
	hasVoted, err := s.repo.UserHasVoted(commentID, userID)
	if err != nil {
		return err
	}
	if hasVoted {
		return nil
	}
	err = s.repo.CreateVote(commentID, userID, false)
	if err != nil {
		return err
	}
	return s.repo.DownvoteComment(commentID)
}
