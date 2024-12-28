package request

// CreateReviewRequest 创建评论请求
type CreateReviewRequest struct {
	BookID  uint   `json:"book_id" binding:"required,min=1"`
	Content string `json:"content" binding:"required,min=1,max=1000"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
}

// UpdateReviewRequest 更新评论请求
type UpdateReviewRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
}

// ReviewSearchRequest 评论搜索请求
type ReviewSearchRequest struct {
	UserID  uint   `form:"user_id" binding:"omitempty,min=1"`
	BookID  uint   `form:"book_id" binding:"omitempty,min=1"`
	Rating  int    `form:"rating" binding:"omitempty,min=1,max=5"`
	Status  *int   `form:"status" binding:"omitempty,oneof=0 1"`
	OrderBy string `form:"order_by" binding:"omitempty,oneof=rating created_at"`
	SearchRequest
}
