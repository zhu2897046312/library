package mysql

import (
	"errors"
	"gorm.io/gorm"
	"library/model"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	List(params *model.SearchParams) ([]*model.User, int64, error)
	Transaction(fc func(tx *gorm.DB) error) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {	
	user.CreatedAt = r.db.NowFunc()
	user.UpdatedAt = r.db.NowFunc()
	user.LastLoginAt = r.db.NowFunc()
	return r.db.Create(user).Error
}

// Update 更新用户信息
func (r *userRepository) Update(user *model.User) error {
	user.UpdatedAt = r.db.NowFunc()
	return r.db.Model(user).Updates(user).Error
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// Transaction wraps the function in a database transaction
func (r *userRepository) Transaction(fc func(tx *gorm.DB) error) error {
	return r.db.Transaction(fc)
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表（支持模糊查询和分页）
func (r *userRepository) List(params *model.SearchParams) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64
	
	db := r.db.Model(&model.User{})
	
	// 模糊查询条件
	if params.Keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
			"%"+params.Keyword+"%",
			"%"+params.Keyword+"%",
			"%"+params.Keyword+"%")
	}
	
	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	err := db.Offset(offset).Limit(params.PageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
}
