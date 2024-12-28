package handler

import (
	"library/handler/request"
	"library/handler/response"
	"library/model"
	"library/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService service.ReviewServiceInterface
}

func NewReviewHandler(reviewService service.ReviewServiceInterface) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

func (h *ReviewHandler) authCheck(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.NewResponse(http.StatusUnauthorized, "用户未登录", nil))
		return 0, false
	}
	return userID.(uint), true
}

// CreateReview 创建评论
// @Summary 创建评论
// @Description 用户创建图书评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body request.CreateReviewRequest true "评论信息"
// @Success 200 {object} response.Response
// @Router /api/v1/reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, ok := h.authCheck(c)
	if !ok {
		return
	}

	var req request.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的请求参数", nil))
		return
	}

	review := &model.Review{
		UserID:  userID,
		BookID:  req.BookID,
		Content: req.Content,
		Rating:  req.Rating,
	}

	if err := h.reviewService.CreateReview( review); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "评论创建成功", review))
}

// UpdateReview 更新评论
// @Summary 更新评论
// @Description 用户更新自己的评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Param request body request.UpdateReviewRequest true "评论信息"
// @Success 200 {object} response.Response
// @Router /api/v1/reviews/{id} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID, ok := h.authCheck(c)
	if !ok {
		return
	}

	var uri request.IDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的评论ID", nil))
		return
	}

	var req request.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的请求参数", nil))
		return
	}

	// 检查是否是评论作者
	review, err := h.reviewService.GetReview( uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	if review.UserID != userID {
		c.JSON(http.StatusForbidden, response.NewResponse(http.StatusForbidden, "没有权限修改此评论", nil))
		return
	}

	review.Content = req.Content
	review.Rating = req.Rating

	if err := h.reviewService.UpdateReview( review); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "评论更新成功", review))
}

// DeleteReview 删除评论
// @Summary 删除评论
// @Description 用户删除自己的评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Router /api/v1/reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID, ok := h.authCheck(c)
	if !ok {
		return
	}

	var uri request.IDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的评论ID", nil))
		return
	}

	// 检查是否是评论作者
	review, err := h.reviewService.GetReview( uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	if review.UserID != userID {
		c.JSON(http.StatusForbidden, response.NewResponse(http.StatusForbidden, "没有权限删除此评论", nil))
		return
	}

	if err := h.reviewService.DeleteReview( uri.ID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "评论删除成功", nil))
}

// GetReview 获取评论详情
// @Summary 获取评论详情
// @Description 获取指定评论的详细信息
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Router /api/v1/reviews/{id} [get]
func (h *ReviewHandler) GetReview(c *gin.Context) {
	var uri request.IDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的评论ID", nil))
		return
	}

	review, err := h.reviewService.GetReview( uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Success", review))
}

// ListReviews 获取评论列表
// @Summary 获取评论列表
// @Description 根据条件搜索评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param request query request.ReviewSearchRequest true "搜索条件"
// @Success 200 {object} response.Response
// @Router /api/v1/reviews [get]
func (h *ReviewHandler) ListReviews(c *gin.Context) {
	var req request.ReviewSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的请求参数", nil))
		return
	}

	searchParams := &model.SearchParams{
		Keyword: req.Keyword,
		OrderBy: req.OrderBy,
	}
	// 设置分页参数
	searchParams.Page = req.Page
	searchParams.PageSize = req.PageSize

	reviews, total, err := h.reviewService.ListReviews( searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Success", gin.H{
		"total": total,
		"items": reviews,
	}))
}

// UpdateReviewStatus 更新评论状态（管理员接口）
// @Summary 更新评论状态
// @Description 管理员更新评论显示状态
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Param request body request.StatusRequest true "状态信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/reviews/{id}/status [put]
func (h *ReviewHandler) UpdateReviewStatus(c *gin.Context) {
	_, ok := h.authCheck(c)
	if !ok {
		return
	}

	role, _ := c.Get("role")
	if role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.NewResponse(http.StatusForbidden, "没有权限修改评论状态", nil))
		return
	}

	var uri request.IDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的评论ID", nil))
		return
	}

	var req request.StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的请求参数", nil))
		return
	}

	review, err := h.reviewService.GetReview( uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	review.Status = req.Status
	if err := h.reviewService.UpdateReview( review); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "评论状态更新成功", nil))
}
