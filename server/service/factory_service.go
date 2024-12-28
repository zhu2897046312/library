package service

import (
	"sync"

	"library/repository/mysql"
)

var (
	factoryInstance *factory
	once            sync.Once
)

// Factory 定义服务工厂接口
type Factory interface {
	GetUserService() UserServiceInterface
	GetReviewService() ReviewServiceInterface
	GetBorrowService() BorrowServiceInterface
	GetBookService() BookServiceInterface
}

// factory 实现Factory接口
type factory struct {
	mysqlFactory mysql.Factory
	userSrv      UserServiceInterface
	reviewSrv    ReviewServiceInterface
	borrowSrv    BorrowServiceInterface
	bookSrv      BookServiceInterface
	mu           sync.RWMutex
}

// NewFactory 创建服务工厂实例（单例))
func NewFactory(mysqlFactory mysql.Factory) Factory {
	once.Do(func() {
		factoryInstance = &factory{
			mysqlFactory: mysqlFactory,
		}
	})
	return factoryInstance
}

func (f *factory) GetUserService() UserServiceInterface {
	f.mu.RLock()
	if f.userSrv != nil {
		defer f.mu.RUnlock()
		return f.userSrv
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.userSrv == nil {
		f.userSrv = NewUserService(f.mysqlFactory.GetUserRepository())
	}
	return f.userSrv
}

func (f *factory) GetReviewService() ReviewServiceInterface {
	f.mu.RLock()
	if f.reviewSrv != nil {
		defer f.mu.RUnlock()
		return f.reviewSrv
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.reviewSrv == nil {
		f.reviewSrv = NewReviewService(
			f.mysqlFactory.GetReviewRepository(),
			f.mysqlFactory.GetBookRepository(),
			f.mysqlFactory.GetUserRepository(),
		)
	}
	return f.reviewSrv
}

func (f *factory) GetBorrowService() BorrowServiceInterface {
	f.mu.RLock()
	if f.borrowSrv != nil {
		defer f.mu.RUnlock()
		return f.borrowSrv
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.borrowSrv == nil {
		f.borrowSrv = NewBorrowService(f.mysqlFactory.GetBorrowRepository(), f.mysqlFactory.GetBookRepository(), f.mysqlFactory.GetUserRepository())
	}
	return f.borrowSrv
}

func (f *factory) GetBookService() BookServiceInterface {
	f.mu.RLock()
	if f.bookSrv != nil {
		defer f.mu.RUnlock()
		return f.bookSrv
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.bookSrv == nil {
		f.bookSrv = NewBookService(
			f.mysqlFactory.GetBookRepository(),
		)
	}
	return f.bookSrv
}
