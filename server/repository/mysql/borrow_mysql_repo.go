package mysql

import (
	"errors"
	"gorm.io/gorm"
	"library/model"
	"time"
)

type BorrowRepository interface {
	Create( borrow *model.Borrow) error
	Update( borrow *model.Borrow) error
	GetByID( id uint) (*model.Borrow, error)
	GetByUserAndBookID( userID, bookID uint) (*model.Borrow, error)
	List( params *model.SearchParams) ([]*model.Borrow, int64, error)
	GetUserBorrows( userID uint, status int) ([]*model.Borrow, error)
	GetOverdueBorrows() ([]*model.Borrow, error)
	Transaction(fc func(tx *gorm.DB) error) error
}


type borrowRepository struct {
	db *gorm.DB
}

func NewBorrowRepository(db *gorm.DB) BorrowRepository {
	return &borrowRepository{db: db}
}

// Transaction wraps the function in a database transaction
func (r *borrowRepository) Transaction(fc func(tx *gorm.DB) error) error {
	return r.db.Transaction(fc)
}

// Create 创建借阅记录
func (r *borrowRepository) Create( borrow *model.Borrow) error {
	return r.db.Create(borrow).Error
}

// Update 更新借阅记录
func (r *borrowRepository) Update( borrow *model.Borrow) error {
	return r.db.Updates(borrow).Error
}

// GetByID 根据ID获取借阅记录
func (r *borrowRepository) GetByID( id uint) (*model.Borrow, error) {
	var borrow model.Borrow
	err := r.db.
		Preload("User").
		Preload("Book").
		First(&borrow, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &borrow, nil
}

// GetByUserAndBookID 根据用户ID和图书ID获取借阅记录
func (r *borrowRepository) GetByUserAndBookID(userID, bookID uint) (*model.Borrow, error) {
	var borrow model.Borrow
	err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&borrow).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &borrow, nil
}

// List 获取借阅记录列表（支持模糊查询和分页）
func (r *borrowRepository) List( params *model.SearchParams) ([]*model.Borrow, int64, error) {
	var borrows []*model.Borrow
	var total int64
	
	db := r.db.Model(&model.Borrow{}).
		Preload("User").
		Preload("Book")
	
	// 时间范围查询
	if params.StartTime != "" {
		db = db.Where("borrow_date >= ?", params.StartTime)
	}
	if params.EndTime != "" {
		db = db.Where("borrow_date <= ?", params.EndTime)
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
	err := db.Offset(offset).Limit(params.PageSize).Find(&borrows).Error
	if err != nil {
		return nil, 0, err
	}
	
	return borrows, total, nil
}

// GetUserBorrows 获取用户的借阅记录
func (r *borrowRepository) GetUserBorrows( userID uint, status int) ([]*model.Borrow, error) {
	var borrows []*model.Borrow
	db := r.db.
		Preload("Book").
		Where("user_id = ?", userID)
	
	if status > 0 {
		db = db.Where("status = ?", status)
	}
	
	err := db.Order("created_at DESC").Find(&borrows).Error
	if err != nil {
		return nil, err
	}
	return borrows, nil
}

// GetOverdueBorrows 获取逾期的借阅记录
func (r *borrowRepository) GetOverdueBorrows() ([]*model.Borrow, error) {
	var borrows []*model.Borrow
	err := r.db.
		Preload("User").
		Preload("Book").
		Where("status = ? AND due_date < ?", 1, time.Now()).
		Find(&borrows).Error
	if err != nil {
		return nil, err
	}
	return borrows, nil
}
