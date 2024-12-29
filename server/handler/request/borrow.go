package request

import "time"

// CreateBorrowRequest 创建借阅请求
type CreateBorrowRequest struct {
	BookID    uint      `json:"book_id" binding:"required,min=1" example:"1"`
	DueDate   time.Time `json:"due_date" binding:"required" example:"2024-01-01T00:00:00+08:00"` // 应还日期必须大于借阅日期
	Remark    string    `json:"remark" binding:"omitempty,max=256" example:"请尽快归还"`
}

// ReturnBookRequest 归还图书请求
type ReturnBookRequest struct {
	BorrowID uint    `json:"borrow_id" binding:"required,min=1" example:"1"`
	Fine     float64 `json:"fine" binding:"omitempty,min=0" example:"10.00"`
	Remark   string  `json:"remark" binding:"omitempty,max=256" example:"请尽快归还"`
}

// UpdateBorrowRequest 更新借阅信息请求
type UpdateBorrowRequest struct {
	DueDate time.Time `json:"due_date" binding:"required" example:"2024-01-01T00:00:00+08:00"`
	Status  int       `json:"status" binding:"required,oneof=0 1 2 3" example:"1"` // 0-已取消 1-借阅中 2-已归还 3-已逾期
	Fine    float64   `json:"fine" binding:"omitempty,min=0" example:"10.00"`
	Remark  string    `json:"remark" binding:"omitempty,max=256" example:"请尽快归还"`
}

// BorrowSearchRequest 借阅记录搜索请求
type BorrowSearchRequest struct {
	UserID     uint      `form:"user_id" binding:"omitempty,min=1" example:"1"`
	BookID     uint      `form:"book_id" binding:"omitempty,min=1" example:"1"`
	Status     []int     `form:"status" binding:"omitempty,dive,oneof=0 1 2 3" example:"1"`
	StartTime  time.Time `form:"start_time" binding:"omitempty" example:"2024-01-01T00:00:00+08:00"`
	EndTime    time.Time `form:"end_time" binding:"omitempty,gtefield=StartTime" example:"2024-01-01T00:00:00+08:00"`
	SearchRequest
}
