package service

import "errors"

var (
	// ErrInvalidParameter 参数错误
	ErrInvalidParameter = errors.New("invalid parameter")
	// ErrNotFound 资源不存在
	ErrNotFound = errors.New("resource not found")
	// ErrAlreadyExists 资源已存在
	ErrAlreadyExists = errors.New("resource already exists")
	// ErrPasswordIncorrect 密码错误
	ErrPasswordIncorrect = errors.New("password incorrect")
	// ErrBookNotAvailable 图书不可借
	ErrBookNotAvailable = errors.New("book not available")
	// ErrBorrowLimitExceeded 超出借阅限制
	ErrBorrowLimitExceeded = errors.New("borrow limit exceeded")
	// ErrNotBorrowed 图书未借出
	ErrNotBorrowed = errors.New("book not borrowed")
	// ErrPermissionDenied 权限不足
	ErrPermissionDenied = errors.New("permission denied")
)
