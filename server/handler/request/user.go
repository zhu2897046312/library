package request

// // RegisterRequest 用户注册请求
// type RegisterRequest struct {
// 	Username string `json:"username" binding:"required,min=3,max=32"`
// 	Password string `json:"password" binding:"required,min=6,max=32"`
// 	Email    string `json:"email" binding:"required,email"`
// 	Phone    string `json:"phone" binding:"omitempty,len=11"`
// 	Nickname string `json:"nickname" binding:"omitempty,min=2,max=32"`
// 	Role     string `json:"role" binding:"required,oneof=user admin"`
// }

// RegisterRequest 用户注册请求
// @Description 用户注册请求参数
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32" example:"zhangsan"`    // 用户名(3-32个字符)
	Password string `json:"password" binding:"required,min=6,max=32" example:"password123"` // 密码(6-32个字符)
	Email    string `json:"email" binding:"required,email" example:"zhangsan@example.com"`  // 邮箱地址
	Phone    string `json:"phone" binding:"omitempty,len=11" example:"13800138000"`         // 手机号(11位)
	Nickname string `json:"nickname" binding:"omitempty,min=2,max=32" example:"张三"`         // 昵称(2-32个字符)
	Role     string `json:"role" binding:"required,oneof=user admin" example:"user"`        // 角色(user/admin)
}

// LoginRequest 用户登录请求
// @Description 用户登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32" example:"zhangsan"`    // 用户名(3-32个字符)
	Password string `json:"password" binding:"required,min=6,max=32" example:"password123"` // 密码(6-32个字符)
}

// UpdateUserRequest 更新用户信息请求
// @Description 更新用户信息请求参数
type UpdateUserRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,min=2,max=32" example:"张三"`         // 昵称(2-32个字符)
	Email    string `json:"email" binding:"omitempty,email" example:"zhangsan@example.com"` // 邮箱地址
	Phone    string `json:"phone" binding:"omitempty,len=11" example:"13800138000"`         // 手机号(11位)
}

// ChangePasswordRequest 修改密码请求
// @Description 修改密码请求参数
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=32" example:"password123"`                        // 旧密码(6-32个字符)
	NewPassword string `json:"new_password" binding:"required,min=6,max=32,nefield=OldPassword" example:"newpassword123"` // 新密码(6-32个字符，不能与旧密码相同)
}

// UpdateUserRoleRequest 更新用户角色请求
// @Description 更新用户角色请求参数
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user admin" example:"user"` // user-普通用户 admin-管理员
}

// UserSearchRequest 用户搜索请求
// @Description 用户搜索请求参数
type UserSearchRequest struct {
	Role    *int   `form:"role" binding:"omitempty,oneof=0 1" example:"0"`   // 0-普通用户 1-管理员
	Status  *int   `form:"status" binding:"omitempty,oneof=0 1" example:"0"` // 0-正常 1-禁用
	Keyword string `form:"keyword" binding:"omitempty,min=1" example:"张三"`   // 关键词
	SearchRequest
}
