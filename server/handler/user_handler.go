package handler

import (
	"library/handler/request"
	"library/handler/response"
	"library/middleware"
	"library/model"
	"library/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	err := h.userService.Register( req.Username, req.Password, req.Email, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "User registered successfully", nil))
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录并获取token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "登录信息"
// @Success 200 {object} response.Response
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	user, err := h.userService.Login( req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.NewResponse(http.StatusUnauthorized, err.Error(), nil))
		return
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, "Failed to generate token", nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Login successful", gin.H{
		"token": token,
		"user":  user,
	}))
}

// GetProfile 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户的访问令牌"
// @Success 200 {object} response.Response
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.NewResponse(http.StatusUnauthorized, "User not authenticated", nil))
		return
	}

	user, err := h.userService.GetUserInfo( userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Success", user))
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户的信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户的访问令牌"
// @Param request body request.UpdateUserRequest true "用户信息"
// @Success 200 {object} response.Response
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	userID, _ := c.Get("userID")
	user, err := h.userService.GetUserInfo( userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := h.userService.UpdateUserInfo( user); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Profile updated successfully", user))
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户的访问令牌"
// @Param request body request.ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response
// @Router /users/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	userID, _ := c.Get("userID")
	if err := h.userService.ChangePassword( userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Password changed successfully", nil))
}

// ListUsers 获取用户列表（管理员接口）
// @Summary 获取用户列表
// @Description 管理员获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户的访问令牌"
// @Param request query request.UserSearchRequest true "搜索条件"
// @Success 200 {object} response.Response
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req request.UserSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	searchParams := &model.SearchParams{
		Keyword:    req.Keyword,
	}

	searchParams.Page = req.Page
	searchParams.PageSize = req.PageSize

	users, total, err := h.userService.ListUsers( searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewPaginationResponse(users, total, req.Page, req.PageSize))
}

// UpdateUserRole 更新用户角色（管理员接口）
// @Summary 更新用户角色
// @Description 管理员更新用户角色 "user" or "admin"
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户的访问令牌"
// @Param id path int true "用户ID"
// @Param request body request.UpdateUserRoleRequest true "角色信息"
// @Success 200 {object} response.Response
// @Router /users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	var uri request.IDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid user ID", nil))
		return
	}

	var req request.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	user, err := h.userService.GetUserInfo( uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	user.Role = req.Role
	if err := h.userService.UpdateUserInfo( user); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "User role updated successfully", user))
}
