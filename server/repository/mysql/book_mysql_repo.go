package mysql

import (
	"errors"
	"gorm.io/gorm"
	"library/model"
)

type BookRepository interface {
	Create( book *model.Book) error
	Update( book *model.Book) error
	Delete( id uint) error
	GetByID( id uint) (*model.Book, error)
	GetByISBN( isbn string) (*model.Book, error)
	List( params *model.SearchParams) ([]*model.Book, int64, error)
	UpdateStock( id uint, available int) error
	Transaction(fc func(tx *gorm.DB) error) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

// Transaction wraps the function in a database transaction
func (r *bookRepository) Transaction(fc func(tx *gorm.DB) error) error {
	return r.db.Transaction(fc)
}

// Create 创建图书
func (r *bookRepository) Create( book *model.Book) error {
	return r.db.Create(book).Error
}

// Update 更新图书信息
func (r *bookRepository) Update( book *model.Book) error {
	return r.db.Updates(book).Error
}

// Delete 删除图书（软删除）
func (r *bookRepository) Delete( id uint) error {
	return r.db.Delete(&model.Book{}, id).Error
}

// GetByID 根据ID获取图书
func (r *bookRepository) GetByID( id uint) (*model.Book, error) {
	var book model.Book
	err := r.db.First(&book, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

// GetByISBN 根据ISBN获取图书
func (r *bookRepository) GetByISBN( isbn string) (*model.Book, error) {
	var book model.Book
	err := r.db.Where("isbn = ?", isbn).First(&book).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

// List 获取图书列表（支持模糊查询和分页）
func (r *bookRepository) List( params *model.SearchParams) ([]*model.Book, int64, error) {
	var books []*model.Book
	var total int64
	
	db := r.db.Model(&model.Book{})
	
	// 模糊查询条件
	if params.Keyword != "" {
		db = db.Where("title LIKE ? OR author LIKE ? OR publisher LIKE ? OR isbn LIKE ?",
			"%"+params.Keyword+"%",
			"%"+params.Keyword+"%",
			"%"+params.Keyword+"%",
			"%"+params.Keyword+"%")
	}
	
	// 分类筛选
	if params.Category != "" {
		db = db.Where("category = ?", params.Category)
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
	}
	
	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	err := db.Offset(offset).Limit(params.PageSize).Find(&books).Error
	if err != nil {
		return nil, 0, err
	}
	
	return books, total, nil
}

// UpdateStock 更新图书库存
func (r *bookRepository) UpdateStock( id uint, available int) error {
	return r.db.Model(&model.Book{}).
		Where("id = ?", id).
		Update("available", available).Error
}
