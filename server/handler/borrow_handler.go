package handler

import (
	"library/handler/request"
	"library/handler/response"
	"library/model"
	"library/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BorrowHandler struct {
	borrowService service.BorrowServiceInterface
}

func NewBorrowHandler(borrowService service.BorrowServiceInterface) *BorrowHandler {
	return &BorrowHandler{
		borrowService: borrowService,
	}
}

func (h *BorrowHandler) authCheck(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.NewResponse(http.StatusUnauthorized, "User not logged in", nil))
		return 0, false
	}
	return userID.(uint), true
}

func (h *BorrowHandler) adminCheck(c *gin.Context) bool {
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.NewResponse(http.StatusForbidden, "Permission denied", nil))
		return false
	}
	return true
}

// BorrowBook 借阅图书
// @Summary 借阅图书
// @Description 用户借阅图书
// @Tags 借阅管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body request.CreateBorrowRequest true "借阅信息"
// @Success 200 {object} response.Response
// @Router /api/v1/borrows [post]
func (h *BorrowHandler) BorrowBook(c *gin.Context) {
	userID, ok := h.authCheck(c)
	if !ok {
		return
	}
	var req request.CreateBorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	if err := h.borrowService.BorrowBook( userID, req.BookID); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Book borrowed successfully", nil))
}

// ReturnBook 归还图书
// @Summary 归还图书
// @Description 用户归还图书
// @Tags 借阅管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body request.ReturnBookRequest true "归还信息"
// @Success 200 {object} response.Response
// @Router /api/v1/borrows/return [post]
func (h *BorrowHandler) ReturnBook(c *gin.Context) {
	userID, ok := h.authCheck(c)
	if !ok {
		return
	}
	var req request.ReturnBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	if err := h.borrowService.ReturnBook( userID, req.BorrowID); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Book returned successfully", nil))
}

// GetBorrow 获取借阅详情
// @Summary 获取借阅详情
// @Description 获取指定借阅记录的详细信息
// @Tags 借阅管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "借阅ID"
// @Success 200 {object} response.Response
// @Router /api/v1/borrows/{id} [get]
func (h *BorrowHandler) GetBorrow(c *gin.Context) {
	userID, ok := h.authCheck(c)
	if !ok {
		return
	}
	var uri request.IDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid borrow ID", nil))
		return
	}

	borrow, err := h.borrowService.GetBorrowInfo( uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	// 检查是否是管理员或借阅者本人
	if !h.adminCheck(c) && userID != borrow.UserID {
		c.JSON(http.StatusForbidden, response.NewResponse(http.StatusForbidden, "Permission denied", nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Success", borrow))
}

// ListBorrows 获取借阅列表
// @Summary 获取借阅列表
// @Description 根据条件搜索借阅记录
// @Tags 借阅管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request query request.BorrowSearchRequest true "搜索条件"
// @Success 200 {object} response.Response
// @Router /api/v1/borrows [get]
func (h *BorrowHandler) ListBorrows(c *gin.Context) {
	userID, ok := h.authCheck(c)
	if !ok {
		return
	}
	var req request.BorrowSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	// 非管理员只能查看自己的借阅记录
	if !h.adminCheck(c) {
		req.UserID = userID
	}

	searchParams := &model.SearchParams{
		Keyword:  req.Keyword,
		StartTime: req.StartTime.Format("2006-01-02"),
		EndTime:   req.EndTime.Format("2006-01-02"),
	}
	// 设置分页参数
	searchParams.Page = req.Page
	searchParams.PageSize = req.PageSize

	borrows, total, err := h.borrowService.ListBorrows( searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Success", gin.H{
		"total": total,
		"items": borrows,
	}))
}

// UpdateBorrow 更新借阅信息（管理员接口）
// @Summary 更新借阅信息
// @Description 管理员更新借阅记录信息
// @Tags 借阅管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "借阅ID"
// @Param request body request.UpdateBorrowRequest true "借阅信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/borrows/{id} [put]
func (h *BorrowHandler) UpdateBorrow(c *gin.Context) {
	if !h.adminCheck(c) {
		return
	}
	var uri request.IDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid borrow ID", nil))
		return
	}

	var req request.UpdateBorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	borrow, err := h.borrowService.GetBorrowInfo( uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	// 更新借阅信息
	borrow.DueDate = req.DueDate
	borrow.Status = req.Status
	borrow.Fine = req.Fine
	borrow.Remark = req.Remark

	if err := h.borrowService.UpdateBorrow( borrow); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Borrow record updated successfully", borrow))
}
