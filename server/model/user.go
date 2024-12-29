package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
// @Description 用户信息
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`                                                                                          // 用户ID
	CreatedAt time.Time      `json:"created_at"`                                                                                                    // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                                                                                                    // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" format:"date-time" example:"2024-01-01T00:00:00+08:00"` // 删除时间

	Username    string    `gorm:"type:varchar(32);uniqueIndex;not null" json:"username"` // 用户名
	Password    string    `gorm:"type:varchar(128);not null" json:"-"`                   // 密码
	Salt        string    `gorm:"type:varchar(128);not null" json:"-"`                   // 密码盐值
	Nickname    string    `gorm:"type:varchar(32)" json:"nickname"`                      // 昵称
	Email       string    `gorm:"type:varchar(128);uniqueIndex" json:"email"`            // 邮箱
	Phone       string    `gorm:"type:varchar(20)" json:"phone"`                         // 手机号
	Role        string    `gorm:"type:varchar(32);default:0;not null" json:"role"`       // 角色 user-普通用户 admin-管理员
	Status      int       `gorm:"type:tinyint;default:1;not null" json:"status"`         // 状态 2-禁用 1-启用
	LastLoginAt time.Time `gorm:"type:datetime" json:"last_login_at"`                    // 最后登录时间
}
