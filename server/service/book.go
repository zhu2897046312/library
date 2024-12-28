package service

import (
	"fmt"
	"gorm.io/gorm"
	"library/model"
	"library/repository/mysql"
)

// BookServiceInterface 图书服务接口
type BookServiceInterface interface {
	CreateBook( book *model.Book) error
	UpdateBook( book *model.Book) error
	DeleteBook( id uint) error
	GetBook( id uint) (*model.Book, error)
	ListBooks( params *model.SearchParams) ([]*model.Book, int64, error)
	UpdateBookStatus( id uint, status int) error
	UpdateBookStock( id uint, change int) error
}


type BookService struct {
	bookRepo mysql.BookRepository
}

func NewBookService(bookRepo mysql.BookRepository) BookServiceInterface {
	return &BookService{
		bookRepo: bookRepo,
	}
}

// CreateBook 创建图书
func (s *BookService) CreateBook( book *model.Book) error {
	return s.bookRepo.Transaction(func(tx *gorm.DB) error {
		// 检查ISBN是否已存在
		existBook, err := s.bookRepo.GetByISBN( book.ISBN)
		if err != nil {
			return fmt.Errorf("check ISBN exists: %w", err)
		}
		if existBook != nil {
			return ErrAlreadyExists
		}

		// 设置初始可借数量
		book.Available = book.Total
		book.Status = 1 // 默认上架

		if err := s.bookRepo.Create( book); err != nil {
			return fmt.Errorf("create book: %w", err)
		}
		return nil
	})
}

// UpdateBook 更新图书信息
func (s *BookService) UpdateBook( book *model.Book) error {
	return s.bookRepo.Transaction(func(tx *gorm.DB) error {
		existBook, err := s.bookRepo.GetByID( book.ID)
		if err != nil {
			return fmt.Errorf("get book by id: %w", err)
		}
		if existBook == nil {
			return ErrNotFound
		}

		// 如果修改了总数量，同步更新可借数量
		if book.Total != existBook.Total {
			diff := book.Total - existBook.Total
			book.Available = existBook.Available + diff
			if book.Available < 0 {
				return fmt.Errorf("invalid total books: would result in negative available books")
			}
		}

		if err := s.bookRepo.Update( book); err != nil {
			return fmt.Errorf("update book: %w", err)
		}
		return nil
	})
}

// DeleteBook 删除图书
func (s *BookService) DeleteBook( id uint) error {
	return s.bookRepo.Transaction(func(tx *gorm.DB) error {
		book, err := s.bookRepo.GetByID( id)
		if err != nil {
			return fmt.Errorf("get book by id: %w", err)
		}
		if book == nil {
			return ErrNotFound
		}

		// 检查是否有未归还的借阅记录
		if book.Available != book.Total {
			return fmt.Errorf("cannot delete book: there are unreturned copies")
		}

		if err := s.bookRepo.Delete( id); err != nil {
			return fmt.Errorf("delete book: %w", err)
		}
		return nil
	})
}

// GetBook 获取图书信息
func (s *BookService) GetBook( id uint) (*model.Book, error) {
	book, err := s.bookRepo.GetByID( id)
	if err != nil {
		return nil, fmt.Errorf("get book by id: %w", err)
	}
	if book == nil {
		return nil, ErrNotFound
	}
	return book, nil
}

// ListBooks 获取图书列表
func (s *BookService) ListBooks( params *model.SearchParams) ([]*model.Book, int64, error) {
	books, total, err := s.bookRepo.List( params)
	if err != nil {
		return nil, 0, fmt.Errorf("list books: %w", err)
	}
	return books, total, nil
}

// UpdateBookStatus 更新图书状态
func (s *BookService) UpdateBookStatus( id uint, status int) error {
	return s.bookRepo.Transaction(func(tx *gorm.DB) error {
		book, err := s.bookRepo.GetByID( id)
		if err != nil {
			return fmt.Errorf("get book by id: %w", err)
		}
		if book == nil {
			return ErrNotFound
		}

		book.Status = status
		if err := s.bookRepo.Update( book); err != nil {
			return fmt.Errorf("update book status: %w", err)
		}
		return nil
	})
}

// UpdateBookStock 更新图书库存
func (s *BookService) UpdateBookStock( id uint, change int) error {
	return s.bookRepo.Transaction(func(tx *gorm.DB) error {
		book, err := s.bookRepo.GetByID( id)
		if err != nil {
			return fmt.Errorf("get book by id: %w", err)
		}
		if book == nil {
			return ErrNotFound
		}

		book.Total += change
		book.Available += change
		if book.Available < 0 || book.Total < 0 {
			return fmt.Errorf("invalid stock change: would result in negative books")
		}

		if err := s.bookRepo.Update( book); err != nil {
			return fmt.Errorf("update book stock: %w", err)
		}
		return nil
	})
}
