package model

import (
	"gorm.io/gorm"
)

// Book 图书模型
type Book struct {
	gorm.Model
	ISBN        string  `gorm:"type:varchar(20);uniqueIndex;not null" json:"isbn"`         // ISBN编号
	Title       string  `gorm:"type:varchar(128);not null" json:"title"`                   // 书名
	Author      string  `gorm:"type:varchar(64);not null" json:"author"`                   // 作者
	Publisher   string  `gorm:"type:varchar(64)" json:"publisher"`                         // 出版社
	Category    string  `gorm:"type:varchar(32)" json:"category"`                          // 分类
	Price       float64 `gorm:"type:decimal(10,2)" json:"price"`                          // 价格
	Total       int     `gorm:"type:int;not null" json:"total"`                           // 总数量
	Available   int     `gorm:"type:int;not null" json:"available"`                       // 可借数量
	Location    string  `gorm:"type:varchar(32)" json:"location"`                         // 馆藏位置
	Cover       string  `gorm:"type:varchar(256)" json:"cover"`                           // 封面图片URL
	Summary     string  `gorm:"type:text" json:"summary"`                                 // 简介
	Status      int     `gorm:"type:tinyint;default:1;not null" json:"status"`           // 状态 0-下架 1-上架
}
