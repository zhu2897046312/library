package model

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`              // 响应码
	Message string      `json:"message"`           // 响应信息
	Data    interface{} `json:"data,omitempty"`    // 响应数据
}

// Pagination 分页参数
type Pagination struct {
	Page     int   `json:"page" form:"page"`         // 页码
	PageSize int   `json:"page_size" form:"page_size"` // 每页数量
	Total    int64 `json:"total"`                    // 总数
}

// SearchParams 通用搜索参数
type SearchParams struct {
	Keyword    string `json:"keyword" form:"keyword"`       // 关键词
	Category   string `json:"category" form:"category"`     // 分类
	StartTime  string `json:"start_time" form:"start_time"` // 开始时间
	EndTime    string `json:"end_time" form:"end_time"`     // 结束时间
	OrderBy    string `json:"order_by" form:"order_by"`     // 排序字段
	OrderType  string `json:"order_type" form:"order_type"` // 排序方式
	Pagination        // 嵌入分页参数
}
