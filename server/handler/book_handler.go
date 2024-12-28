package handler

import (
	"library/handler/request"
	"library/handler/response"
	"library/model"
	"library/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookService service.BookServiceInterface
}

func NewBookHandler(bookService service.BookServiceInterface) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) authMiddleware(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.NewResponse(http.StatusForbidden, "Permission denied", nil))
		c.Abort()
		return
	}
}

// CreateBook 创建图书
// @Summary 创建新图书
// @Description 管理员创建新图书
// @Tags 图书管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body request.CreateBookRequest true "图书信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	h.authMiddleware(c)
	var req request.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	book := &model.Book{
		ISBN:      req.ISBN,
		Title:     req.Title,
		Author:    req.Author,
		Publisher: req.Publisher,
		Category:  req.Category,
		Price:     req.Price,
		Total:     req.Total,
		Location:  req.Location,
		Cover:     req.Cover,
		Summary:   req.Summary,
		Status:    1, // 默认上架
	}

	if err := h.bookService.CreateBook( book); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Book created successfully", book))
}

// UpdateBook 更新图书信息
// @Summary 更新图书信息
// @Description 管理员更新图书信息
// @Tags 图书管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "图书ID"
// @Param request body request.UpdateBookRequest true "图书信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	h.authMiddleware(c)
	var uri request.IDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid book ID", nil))
		return
	}

	var req request.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	book, err := h.bookService.GetBook( uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	// 更新图书信息
	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.Publisher != "" {
		book.Publisher = req.Publisher
	}
	if req.Category != "" {
		book.Category = req.Category
	}
	if req.Price > 0 {
		book.Price = req.Price
	}
	if req.Total > 0 {
		book.Total = req.Total
	}
	if req.Location != "" {
		book.Location = req.Location
	}
	if req.Cover != "" {
		book.Cover = req.Cover
	}
	if req.Summary != "" {
		book.Summary = req.Summary
	}

	if err := h.bookService.UpdateBook( book); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Book updated successfully", book))
}

// GetBook 获取图书详情
// @Summary 获取图书详情
// @Description 获取指定图书的详细信息
// @Tags 图书管理
// @Accept json
// @Produce json
// @Param id path int true "图书ID"
// @Success 200 {object} response.Response
// @Router /api/v1/books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
	var uri request.IDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid book ID", nil))
		return
	}

	book, err := h.bookService.GetBook( uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Success", book))
}

// ListBooks 获取图书列表
// @Summary 获取图书列表
// @Description 根据条件搜索图书
// @Tags 图书管理
// @Accept json
// @Produce json
// @Param request query request.BookSearchRequest true "搜索条件"
// @Success 200 {object} response.Response
// @Router /api/v1/books [get]
func (h *BookHandler) ListBooks(c *gin.Context) {
	var req request.BookSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	searchParams := &model.SearchParams{
		Keyword:  req.Keyword,
		Category: req.Category,
	}
	// 设置分页参数
	searchParams.Page = req.Page
	searchParams.PageSize = req.PageSize

	books, total, err := h.bookService.ListBooks( searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Success", gin.H{
		"total": total,
		"items": books,
	}))
}

// UpdateBookStatus 更新图书状态
// @Summary 更新图书状态
// @Description 管理员更新图书上下架状态
// @Tags 图书管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "图书ID"
// @Param request body request.StatusRequest true "状态信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/books/{id}/status [put]
func (h *BookHandler) UpdateBookStatus(c *gin.Context) {
	h.authMiddleware(c)
	var uri request.IDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid book ID", nil))
		return
	}

	var req request.StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	if err := h.bookService.UpdateBookStatus( uri.ID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Book status updated successfully", nil))
}

// UpdateBookStock 更新图书库存
// @Summary 更新图书库存
// @Description 管理员更新图书库存数量
// @Tags 图书管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "图书ID"
// @Param request body request.UpdateBookStockRequest true "库存信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/books/{id}/stock [put]
func (h *BookHandler) UpdateBookStock(c *gin.Context) {
	h.authMiddleware(c)
	var uri request.IDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid book ID", nil))
		return
	}

	var req request.UpdateBookStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "Invalid request parameters", nil))
		return
	}

	if err := h.bookService.UpdateBookStock( uri.ID, req.Change); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Book stock updated successfully", nil))
}