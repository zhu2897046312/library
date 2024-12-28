package mysql

import (
	"sync"

	"gorm.io/gorm"
)

var (
	factoryInstance *factory
	once           sync.Once
)

// Factory 定义仓库工厂接口
type Factory interface {
	GetUserRepository() UserRepository
	GetReviewRepository() ReviewRepository
	GetBorrowRepository() BorrowRepository
	GetBookRepository() BookRepository
}

// factory 实现Factory接口
type factory struct {
	db          *gorm.DB
	userRepo    UserRepository
	reviewRepo  ReviewRepository
	borrowRepo  BorrowRepository
	bookRepo    BookRepository
	mu          sync.RWMutex
}

// NewFactory 创建工厂实例（单例))
func NewFactory(db *gorm.DB) Factory {
	once.Do(func() {
		factoryInstance = &factory{
			db: db,
		}
	})
	return factoryInstance
}

func (f *factory) GetUserRepository() UserRepository {
	f.mu.RLock()
	if f.userRepo != nil {
		defer f.mu.RUnlock()
		return f.userRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.userRepo == nil {
		f.userRepo = NewUserRepository(f.db)
	}
	return f.userRepo
}

func (f *factory) GetReviewRepository() ReviewRepository {
	f.mu.RLock()
	if f.reviewRepo != nil {
		defer f.mu.RUnlock()
		return f.reviewRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.reviewRepo == nil {
		f.reviewRepo = NewReviewRepository(f.db)
	}
	return f.reviewRepo
}

func (f *factory) GetBorrowRepository() BorrowRepository {
	f.mu.RLock()
	if f.borrowRepo != nil {
		defer f.mu.RUnlock()
		return f.borrowRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.borrowRepo == nil {
		f.borrowRepo = NewBorrowRepository(f.db)
	}
	return f.borrowRepo
}

func (f *factory) GetBookRepository() BookRepository {
	f.mu.RLock()
	if f.bookRepo != nil {
		defer f.mu.RUnlock()
		return f.bookRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.bookRepo == nil {
		f.bookRepo = NewBookRepository(f.db)
	}
	return f.bookRepo
}

