package request

// CreateBookRequest 创建图书请求
// @Description 创建新图书的请求参数
type CreateBookRequest struct {
	ISBN      string  `json:"isbn" binding:"required,min=10,max=13" example:"9787111111111"`
	Title     string  `json:"title" binding:"required,min=1,max=128" example:"The Catcher in the Rye"`
	Author    string  `json:"author" binding:"required,min=1,max=64" example:"J.D. Salinger"`
	Publisher string  `json:"publisher" binding:"required,min=1,max=64" example:"Little, Brown and Company"`
	Category  string  `json:"category" binding:"required,min=1,max=32" example:"Fiction"`
	Price     float64 `json:"price" binding:"required,min=0" example:"10.00"`
	Total     int     `json:"total" binding:"required,min=1" example:"100"`
	Location  string  `json:"location" binding:"required,min=1,max=32" example:"Shelf A1"`
	Cover     string  `json:"cover" binding:"omitempty,url" example:"https://example.com/cover.jpg"`
	Summary   string  `json:"summary" binding:"omitempty,max=1000" example:"This is a great book about life and love."`
}

// UpdateBookRequest 更新图书请求
// @Description 更新图书信息的请求参数
type UpdateBookRequest struct {
	Title     string  `json:"title" binding:"omitempty,min=1,max=128" example:"The Catcher in the Rye"`
	Author    string  `json:"author" binding:"omitempty,min=1,max=64" example:"J.D. Salinger"`
	Publisher string  `json:"publisher" binding:"omitempty,min=1,max=64" example:"Little, Brown and Company"`
	Category  string  `json:"category" binding:"omitempty,min=1,max=32" example:"Fiction"`
	Price     float64 `json:"price" binding:"omitempty,min=0" example:"10.00"`
	Total     int     `json:"total" binding:"omitempty,min=0" example:"100"`
	Location  string  `json:"location" binding:"omitempty,min=1,max=32" example:"Shelf A1"`
	Cover     string  `json:"cover" binding:"omitempty,url" example:"https://example.com/cover.jpg"`
	Summary   string  `json:"summary" binding:"omitempty,max=1000" example:"This is a great book about life and love."`
}

// UpdateBookStockRequest 更新图书库存请求
// @Description 更新图书库存的请求参数
type UpdateBookStockRequest struct {
	Change int `json:"change" binding:"required" example:"10"` // 可以为负数，表示减少库存
}

// BookSearchRequest 图书搜索请求
// @Description 搜索图书的请求参数
type BookSearchRequest struct {
	Category  string  `form:"category" binding:"omitempty,min=1,max=32" example:"Fiction"`
	MinPrice  float64 `form:"min_price" binding:"omitempty,min=0" example:"10.00"`
	MaxPrice  float64 `form:"max_price" binding:"omitempty,min=0,gtefield=MinPrice" example:"20.00"`
	Available *bool   `form:"available" binding:"omitempty" example:"true"` // true: 只显示可借阅的图书
	Status    *int    `form:"status" binding:"omitempty,oneof=0 1" example:"0"`
	SearchRequest
}
