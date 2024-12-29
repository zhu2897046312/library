package model

import (
	"time"

	"gorm.io/gorm"
)

// Review 图书评论模型
// @Description 评论信息
type Review struct {
	ID        uint           `gorm:"primarykey" json:"id"`                                                                                          // 评论ID
	CreatedAt time.Time      `json:"created_at"`                                                                                                    // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                                                                                                    // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" format:"date-time" example:"2024-01-01T00:00:00+08:00"` // 删除时间

	UserID  uint   `gorm:"not null;index" json:"user_id"`                 // 用户ID
	BookID  uint   `gorm:"not null;index" json:"book_id"`                 // 图书ID
	Content string `gorm:"type:text" json:"content"`                      // 评论内容
	Rating  int    `gorm:"type:tinyint;not null" json:"rating"`           // 评分(1-5)
	Status  int    `gorm:"type:tinyint;default:1;not null" json:"status"` // 状态 2-隐藏 1-显示

	User User `gorm:"foreignKey:UserID" json:"user"` // 用户信息
	Book Book `gorm:"foreignKey:BookID" json:"book"` // 图书信息
}
