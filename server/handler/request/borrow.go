package request

import "time"

// CreateBorrowRequest 创建借阅请求
type CreateBorrowRequest struct {
	BookID    uint      `json:"book_id" binding:"required,min=1"`
	DueDate   time.Time `json:"due_date" binding:"required,gtfield=BorrowDate"` // 应还日期必须大于借阅日期
	Remark    string    `json:"remark" binding:"omitempty,max=256"`
}

// ReturnBookRequest 归还图书请求
type ReturnBookRequest struct {
	BorrowID uint    `json:"borrow_id" binding:"required,min=1"`
	Fine     float64 `json:"fine" binding:"omitempty,min=0"`
	Remark   string  `json:"remark" binding:"omitempty,max=256"`
}

// UpdateBorrowRequest 更新借阅信息请求
type UpdateBorrowRequest struct {
	DueDate time.Time `json:"due_date" binding:"required"`
	Status  int       `json:"status" binding:"required,oneof=0 1 2 3"` // 0-已取消 1-借阅中 2-已归还 3-已逾期
	Fine    float64   `json:"fine" binding:"omitempty,min=0"`
	Remark  string    `json:"remark" binding:"omitempty,max=256"`
}

// BorrowSearchRequest 借阅记录搜索请求
type BorrowSearchRequest struct {
	UserID     uint      `form:"user_id" binding:"omitempty,min=1"`
	BookID     uint      `form:"book_id" binding:"omitempty,min=1"`
	Status     []int     `form:"status" binding:"omitempty,dive,oneof=0 1 2 3"`
	StartTime  time.Time `form:"start_time" binding:"omitempty"`
	EndTime    time.Time `form:"end_time" binding:"omitempty,gtefield=StartTime"`
	SearchRequest
}
