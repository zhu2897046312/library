package response

// Response represents common API response
type Response struct {
	Code int         `json:"code"`    // HTTP状态码
	Msg  string      `json:"msg"`     // 响应消息
	Data interface{} `json:"data"`    // 响应数据
}

// PaginationData represents pagination information
type PaginationData struct {
	Total    int64       `json:"total"`     // 总数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页数量
	Items    interface{} `json:"items"`     // 数据列表
}

// NewResponse creates a new response
func NewResponse(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// NewPaginationResponse creates a new pagination response
func NewPaginationResponse(items interface{}, total int64, page, pageSize int) *Response {
	return NewResponse(200, "success", PaginationData{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    items,
	})
}
