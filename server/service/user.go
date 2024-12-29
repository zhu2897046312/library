package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gorm.io/gorm"
	"library/model"
	"library/repository/mysql"
	"time"
)

type UserServiceInterface interface {
	Register(username, password, email string , role string) error
	Login(username, password string) (*model.User, error)
	GetUserInfo(id uint) (*model.User, error)
	UpdateUserInfo(user *model.User) error
	ChangePassword(id uint, oldPassword, newPassword string) error
	ListUsers(params *model.SearchParams) ([]*model.User, int64, error)
}


type userService struct {
	userRepo mysql.UserRepository
}

func NewUserService(userRepo mysql.UserRepository) UserServiceInterface {
	return &userService{
		userRepo: userRepo,
	}
}

// Register 用户注册
func (s *userService) Register(username, password, email string, role string) error {
	return s.userRepo.Transaction(func(tx *gorm.DB) error {
		// 检查用户名是否已存在
		existUser, err := s.userRepo.GetByUsername( username)
		if err != nil {
			return fmt.Errorf("check username exists: %w", err)
		}
		if existUser != nil {
			return ErrAlreadyExists
		}

		// 生成密码盐值和加密密码
		salt := generateSalt()
		encryptedPass := encryptPassword(password, salt)

		user := &model.User{
			Username: username,
			Password: encryptedPass,
			Salt:     salt,
			Email:    email,
			Role:     role, // 默认普通用户
			Status:   1, // 默认启用
		}

		if err := s.userRepo.Create( user); err != nil {
			return fmt.Errorf("create user: %w", err)
		}
		return nil
	})
}

// Login 用户登录
func (s *userService) Login(username, password string) (*model.User, error) {
	var user *model.User
	err := s.userRepo.Transaction(func(tx *gorm.DB) error {
		var err error
		user, err = s.userRepo.GetByUsername( username)
		if err != nil {
			return fmt.Errorf("get user by username: %w", err)
		}
		if user == nil {
			return ErrNotFound
		}

		// 验证密码
		if encryptPassword(password, user.Salt) != user.Password {
			return ErrPasswordIncorrect
		}

		// 更新最后登录时间
		user.LastLoginAt = time.Now()
		if err := s.userRepo.Update(user); err != nil {
			return fmt.Errorf("update last login time: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserInfo 获取用户信息
func (s *userService) GetUserInfo(id uint) (*model.User, error) {
	return s.userRepo.GetByID( id)
}

// UpdateUserInfo 更新用户信息
func (s *userService) UpdateUserInfo(user *model.User) error {
	existUser, err := s.userRepo.GetByID( user.ID)
	if err != nil {
		return err
	}
	if existUser == nil {
		return ErrNotFound
	}

	// 保持原有的敏感信息不变
	user.Password = existUser.Password
	user.Salt = existUser.Salt
	user.Role = existUser.Role

	return s.userRepo.Update( user)
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID( id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNotFound
	}

	// 验证旧密码
	if encryptPassword(oldPassword, user.Salt) != user.Password {
		return ErrPasswordIncorrect
	}

	// 更新密码
	user.Salt = generateSalt()
	user.Password = encryptPassword(newPassword, user.Salt)

	return s.userRepo.Update( user)
}

// ListUsers 获取用户列表
func (s *userService) ListUsers(params *model.SearchParams) ([]*model.User, int64, error) {
	return s.userRepo.List( params)
}

// 生成随机盐值
func generateSalt() string {
	// 这里简单使用时间戳作为盐值，实际应用中应使用更安全的随机数
	return hex.EncodeToString([]byte(time.Now().String()))
}

// 加密密码
func encryptPassword(password, salt string) string {
	hash := md5.New()
	hash.Write([]byte(password + salt))
	return hex.EncodeToString(hash.Sum(nil))
}
