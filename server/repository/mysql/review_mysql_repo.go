package mysql

import (
	"errors"
	"gorm.io/gorm"
	"library/model"
)

type ReviewRepository interface {
	Create(review *model.Review) error
	Update(review *model.Review) error
	Delete(id uint) error
	GetByID(id uint) (*model.Review, error)
	List(params *model.SearchParams) ([]*model.Review, int64, error)
	GetBookReviews(bookID uint, params *model.SearchParams) ([]*model.Review, int64, error)
	GetUserReviews(userID uint, params *model.SearchParams) ([]*model.Review, int64, error)
	Transaction(fc func(tx *gorm.DB) error) error
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

// Transaction wraps the function in a database transaction
func (r *reviewRepository) Transaction(fc func(tx *gorm.DB) error) error {
	return r.db.Transaction(fc)
}

// Create 创建评论
func (r *reviewRepository) Create(review *model.Review) error {
	return r.db.Create(review).Error
}

// Update 更新评论
func (r *reviewRepository) Update(review *model.Review) error {
	return r.db.Model(review).Updates(review).Error
}

// Delete 删除评论（软删除）
func (r *reviewRepository) Delete(id uint) error {
	return r.db.Delete(&model.Review{}, id).Error
}

// GetByID 根据ID获取评论
func (r *reviewRepository) GetByID(id uint) (*model.Review, error) {
	var review model.Review
	err := r.db.
		Preload("User").
		Preload("Book").
		First(&review, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

// List 获取评论列表（支持模糊查询和分页）
func (r *reviewRepository) List(params *model.SearchParams) ([]*model.Review, int64, error) {
	var reviews []*model.Review
	var total int64

	db := r.db.Model(&model.Review{}).
		Preload("User").
		Preload("Book")

	// 模糊查询条件
	if params.Keyword != "" {
		db = db.Where("content LIKE ?", "%"+params.Keyword+"%")
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	if params.OrderBy != "" {
		order := params.OrderBy
		if params.OrderType != "" {
			order += " " + params.OrderType
		}
		db = db.Order(order)
	} else {
		db = db.Order("created_at DESC")
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	err := db.Offset(offset).Limit(params.PageSize).Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

// GetBookReviews 获取图书的评论列表
func (r *reviewRepository) GetBookReviews(bookID uint, params *model.SearchParams) ([]*model.Review, int64, error) {
	var reviews []*model.Review
	var total int64

	db := r.db.Model(&model.Review{}).
		Preload("User").
		Where("book_id = ?", bookID)

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	err := db.Offset(offset).
		Limit(params.PageSize).
		Order("created_at DESC").
		Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

// GetUserReviews 获取用户的评论列表
func (r *reviewRepository) GetUserReviews(userID uint, params *model.SearchParams) ([]*model.Review, int64, error) {
	var reviews []*model.Review
	var total int64

	db := r.db.Model(&model.Review{}).
		Preload("Book").
		Where("user_id = ?", userID)

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	err := db.Offset(offset).
		Limit(params.PageSize).
		Order("created_at DESC").
		Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}
