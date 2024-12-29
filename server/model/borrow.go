package model

import (
	"time"

	"gorm.io/gorm"
)

// Borrow 借阅记录模型
// @Description 借阅信息
type Borrow struct {
	ID        uint           `gorm:"primarykey" json:"id"`                                                                                          // 借阅记录ID
	CreatedAt time.Time      `json:"created_at"`                                                                                                    // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                                                                                                    // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" format:"date-time" example:"2024-01-01T00:00:00+08:00"` // 删除时间

	UserID     uint      `gorm:"not null;index" json:"user_id"`                 // 用户ID
	BookID     uint      `gorm:"not null;index" json:"book_id"`                 // 图书ID
	BorrowDate time.Time `gorm:"type:datetime;not null" json:"borrow_date"`     // 借出时间
	DueDate    time.Time `gorm:"type:datetime;not null" json:"due_date"`        // 应还时间
	ReturnDate time.Time `gorm:"type:datetime" json:"return_date"`              // 实际归还时间
	Status     int       `gorm:"type:tinyint;default:1;not null" json:"status"` // 状态 4-已取消 1-借阅中 2-已归还 3-已逾期
	Fine       float64   `gorm:"type:decimal(10,2);default:0" json:"fine"`      // 罚金
	Remark     string    `gorm:"type:varchar(256)" json:"remark"`               // 备注

	User User `gorm:"foreignKey:UserID" json:"user"` // 用户信息
	Book Book `gorm:"foreignKey:BookID" json:"book"` // 图书信息
}
