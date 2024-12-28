package model

import (
	"gorm.io/gorm"
)

// Review 图书评论模型
type Review struct {
	gorm.Model
	UserID    uint    `gorm:"not null;index" json:"user_id"`                    // 用户ID
	BookID    uint    `gorm:"not null;index" json:"book_id"`                    // 图书ID
	Content   string  `gorm:"type:text" json:"content"`                         // 评论内容
	Rating    int     `gorm:"type:tinyint;not null" json:"rating"`             // 评分(1-5)
	Status    int     `gorm:"type:tinyint;default:1;not null" json:"status"`   // 状态 0-隐藏 1-显示
	
	User      User    `gorm:"foreignKey:UserID" json:"user"`                   // 用户信息
	Book      Book    `gorm:"foreignKey:BookID" json:"book"`                   // 图书信息
}
