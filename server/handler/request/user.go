package request

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
	Nickname string `json:"nickname" binding:"omitempty,min=2,max=32"`
	Role     string `json:"role" binding:"required,oneof=user admin"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,min=2,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=32"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32,nefield=OldPassword"`
}

// UpdateUserRoleRequest 更新用户角色请求
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user admin"` // user-普通用户 admin-管理员
}

// UserSearchRequest 用户搜索请求
type UserSearchRequest struct {
	Role    *int   `form:"role" binding:"omitempty,oneof=0 1"`
	Status  *int   `form:"status" binding:"omitempty,oneof=0 1"`
	Keyword string `form:"keyword" binding:"omitempty,min=1"`
	SearchRequest
}
