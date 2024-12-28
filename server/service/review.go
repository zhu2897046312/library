package service

import (
	"library/model"
	"library/repository/mysql"
)

// ReviewServiceInterface 评论服务接口
type ReviewServiceInterface interface {
	CreateReview(review *model.Review) error
	UpdateReview(review *model.Review) error
	DeleteReview(id uint, userID uint) error
	GetReview(id uint) (*model.Review, error)
	ListReviews(params *model.SearchParams) ([]*model.Review, int64, error)
	GetBookReviews(bookID uint, params *model.SearchParams) ([]*model.Review, int64, error)
	GetUserReviews(userID uint, params *model.SearchParams) ([]*model.Review, int64, error)
}


type ReviewService struct {
	reviewRepo mysql.ReviewRepository
	bookRepo   mysql.BookRepository
	userRepo   mysql.UserRepository
}

func NewReviewService(reviewRepo mysql.ReviewRepository, bookRepo mysql.BookRepository, userRepo mysql.UserRepository) ReviewServiceInterface {
	return &ReviewService{
		reviewRepo: reviewRepo,
		bookRepo:   bookRepo,	
		userRepo:   userRepo,
	}
}

// CreateReview 创建评论
func (s *ReviewService) CreateReview(review *model.Review) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID( review.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNotFound
	}

	// 检查图书是否存在
	book, err := s.bookRepo.GetByID( review.BookID)
	if err != nil {
		return err
	}
	if book == nil {
		return ErrNotFound
	}

	// 验证评分范围
	if review.Rating < 1 || review.Rating > 5 {
		return ErrInvalidParameter
	}

	review.Status = 1 // 默认显示
	return s.reviewRepo.Create( review)
}

// UpdateReview 更新评论
func (s *ReviewService) UpdateReview(review *model.Review) error {
	existReview, err := s.reviewRepo.GetByID( review.ID)
	if err != nil {
		return err
	}
	if existReview == nil {
		return ErrNotFound
	}

	// 只允许更新自己的评论
	if existReview.UserID != review.UserID {
		return ErrPermissionDenied
	}

	// 验证评分范围
	if review.Rating < 1 || review.Rating > 5 {
		return ErrInvalidParameter
	}

	return s.reviewRepo.Update( review)
}

// DeleteReview 删除评论
func (s *ReviewService) DeleteReview(id uint, userID uint) error {
	review, err := s.reviewRepo.GetByID( id)
	if err != nil {
		return err
	}
	if review == nil {
		return ErrNotFound
	}

	// 只允许删除自己的评论
	if review.UserID != userID {
		return ErrPermissionDenied
	}

	return s.reviewRepo.Delete( id)
}

// GetReview 获取评论
func (s *ReviewService) GetReview(id uint) (*model.Review, error) {
	return s.reviewRepo.GetByID( id)
}

// ListReviews 获取评论列表
func (s *ReviewService) ListReviews(params *model.SearchParams) ([]*model.Review, int64, error) {
	return s.reviewRepo.List( params)
}

// GetBookReviews 获取图书的评论列表
func (s *ReviewService) GetBookReviews(bookID uint, params *model.SearchParams) ([]*model.Review, int64, error) {
	return s.reviewRepo.GetBookReviews( bookID, params)
}

// GetUserReviews 获取用户的评论列表
func (s *ReviewService) GetUserReviews(userID uint, params *model.SearchParams) ([]*model.Review, int64, error) {
	return s.reviewRepo.GetUserReviews( userID, params)
}
