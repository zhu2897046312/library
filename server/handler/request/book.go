package request

// CreateBookRequest 创建图书请求
type CreateBookRequest struct {
	ISBN      string  `json:"isbn" binding:"required,min=10,max=13"`
	Title     string  `json:"title" binding:"required,min=1,max=128"`
	Author    string  `json:"author" binding:"required,min=1,max=64"`
	Publisher string  `json:"publisher" binding:"required,min=1,max=64"`
	Category  string  `json:"category" binding:"required,min=1,max=32"`
	Price     float64 `json:"price" binding:"required,min=0"`
	Total     int     `json:"total" binding:"required,min=1"`
	Location  string  `json:"location" binding:"required,min=1,max=32"`
	Cover     string  `json:"cover" binding:"omitempty,url"`
	Summary   string  `json:"summary" binding:"omitempty,max=1000"`
}

// UpdateBookRequest 更新图书请求
type UpdateBookRequest struct {
	Title     string  `json:"title" binding:"omitempty,min=1,max=128"`
	Author    string  `json:"author" binding:"omitempty,min=1,max=64"`
	Publisher string  `json:"publisher" binding:"omitempty,min=1,max=64"`
	Category  string  `json:"category" binding:"omitempty,min=1,max=32"`
	Price     float64 `json:"price" binding:"omitempty,min=0"`
	Total     int     `json:"total" binding:"omitempty,min=0"`
	Location  string  `json:"location" binding:"omitempty,min=1,max=32"`
	Cover     string  `json:"cover" binding:"omitempty,url"`
	Summary   string  `json:"summary" binding:"omitempty,max=1000"`
}

// UpdateBookStockRequest 更新图书库存请求
type UpdateBookStockRequest struct {
	Change int `json:"change" binding:"required"` // 可以为负数，表示减少库存
}

// BookSearchRequest 图书搜索请求
type BookSearchRequest struct {
	Category  string  `form:"category" binding:"omitempty,min=1,max=32"`
	MinPrice  float64 `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice  float64 `form:"max_price" binding:"omitempty,min=0,gtefield=MinPrice"`
	Available *bool   `form:"available" binding:"omitempty"` // true: 只显示可借阅的图书
	Status    *int    `form:"status" binding:"omitempty,oneof=0 1"`
	SearchRequest
}
