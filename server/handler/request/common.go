package request

// PaginationRequest 分页请求的通用结构
type PaginationRequest struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}

// SearchRequest 搜索请求的通用结构
type SearchRequest struct {
	Keyword string `form:"keyword" binding:"omitempty,min=1"`
	PaginationRequest
}

// IDRequest 通过ID查询的通用结构
type IDRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// StatusRequest 状态更新的通用结构
type StatusRequest struct {
	Status int `json:"status" binding:"required,oneof=0 1"` // 0:禁用 1:启用
}

// SortRequest 排序的通用结构
type SortRequest struct {
	OrderBy string `form:"order_by" binding:"omitempty,oneof=created_at updated_at"`
	Order   string `form:"order" binding:"omitempty,oneof=asc desc"`
}
