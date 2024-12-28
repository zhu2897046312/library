package service

import (
	"library/model"
	"library/repository/mysql"
	"time"
)

// BorrowServiceInterface 借阅服务接口
type BorrowServiceInterface interface {
	BorrowBook(userID, bookID uint) error
	ReturnBook(userID, bookID uint) error
	RenewBook(id uint) error
	GetBorrow(id uint) (*model.Borrow, error)
	GetBorrowInfo(id uint) (*model.Borrow, error)
	ListBorrows(params *model.SearchParams) ([]*model.Borrow, int64, error)
	GetUserBorrows(userID uint, status int) ([]*model.Borrow, error)
	GetOverdueBorrows() ([]*model.Borrow, error)
	UpdateBorrow(borrow *model.Borrow) error
}

type BorrowService struct {
	borrowRepo mysql.BorrowRepository
	bookRepo   mysql.BookRepository
	userRepo   mysql.UserRepository
}

func NewBorrowService(borrowRepo mysql.BorrowRepository, bookRepo mysql.BookRepository, userRepo mysql.UserRepository) BorrowServiceInterface {
	return &BorrowService{
		borrowRepo: borrowRepo,
		bookRepo:   bookRepo,
		userRepo:   userRepo,
	}
}

// BorrowBook 借阅图书
func (s *BorrowService) BorrowBook(userID, bookID uint) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNotFound
	}

	// 检查图书是否存在且可借
	book, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}
	if book == nil {
		return ErrNotFound
	}
	if book.Status != 1 {
		return ErrBookNotAvailable
	}
	if book.Available <= 0 {
		return ErrBookNotAvailable
	}

	// 检查用户是否有未归还的同一本书
	borrows, err := s.borrowRepo.GetUserBorrows(userID, 1)
	if err != nil {
		return err
	}
	for _, b := range borrows {
		if b.BookID == bookID {
			return ErrAlreadyExists
		}
	}

	// 创建借阅记录
	borrow := &model.Borrow{
		UserID:     userID,
		BookID:     bookID,
		BorrowDate: time.Now(),
		DueDate:    time.Now().AddDate(0, 0, 30), // 默认借期30天
		Status:     1,                            // 借阅中
	}

	// 更新图书库存
	if err := s.bookRepo.UpdateStock(bookID, book.Available-1); err != nil {
		return err
	}

	return s.borrowRepo.Create(borrow)
}

// ReturnBook 归还图书
func (s *BorrowService) ReturnBook(userID, bookID uint) error {
	// 获取借阅记录
	borrow, err := s.borrowRepo.GetByUserAndBookID(userID, bookID)
	if err != nil {
		return err
	}
	if borrow == nil {
		return ErrNotFound
	}
	if borrow.Status != 1 {
		return ErrNotBorrowed
	}

	// 更新借阅状态
	borrow.Status = 2 // 已归还
	borrow.ReturnDate = time.Now()

	// 计算是否逾期及罚金
	if borrow.ReturnDate.After(borrow.DueDate) {
		days := int(borrow.ReturnDate.Sub(borrow.DueDate).Hours() / 24)
		borrow.Fine = float64(days) * 0.5 // 每天罚款0.5元
	}

	// 更新图书库存
	book, err := s.bookRepo.GetByID(borrow.BookID)
	if err != nil {
		return err
	}
	if err := s.bookRepo.UpdateStock(borrow.BookID, book.Available+1); err != nil {
		return err
	}

	return s.borrowRepo.Update(borrow)
}

// RenewBook 续借图书
func (s *BorrowService) RenewBook(borrowID uint) error {
	borrow, err := s.borrowRepo.GetByID(borrowID)
	if err != nil {
		return err
	}
	if borrow == nil {
		return ErrNotFound
	}
	if borrow.Status != 1 {
		return ErrNotBorrowed
	}

	// 更新到期时间（从当前时间起再借30天）
	borrow.DueDate = time.Now().AddDate(0, 0, 30)
	return s.borrowRepo.Update(borrow)
}

// GetBorrow 获取借阅记录
func (s *BorrowService) GetBorrow(id uint) (*model.Borrow, error) {
	return s.borrowRepo.GetByID(id)
}

// GetBorrowInfo 获取借阅信息
func (s *BorrowService) GetBorrowInfo(id uint) (*model.Borrow, error) {
	return s.GetBorrow(id)
}

// ListBorrows 获取借阅记录列表
func (s *BorrowService) ListBorrows(params *model.SearchParams) ([]*model.Borrow, int64, error) {
	return s.borrowRepo.List(params)
}

// GetUserBorrows 获取用户的借阅记录
func (s *BorrowService) GetUserBorrows(userID uint, status int) ([]*model.Borrow, error) {
	return s.borrowRepo.GetUserBorrows(userID, status)
}

// GetOverdueBorrows 获取逾期的借阅记录
func (s *BorrowService) GetOverdueBorrows() ([]*model.Borrow, error) {
	return s.borrowRepo.GetOverdueBorrows()
}

// UpdateBorrow 更新借阅记录
func (s *BorrowService) UpdateBorrow(borrow *model.Borrow) error {
	return s.borrowRepo.Update(borrow)
}
