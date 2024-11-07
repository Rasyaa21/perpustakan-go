package rest

import (
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/server"
	"base-gin/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BorrowHandler struct {
	hr      *server.Handler
	service *service.BorrowService
}

func NewBorrowHandler(handler *server.Handler, borrowService *service.BorrowService) *BorrowHandler {
	return &BorrowHandler{hr: handler, service: borrowService}
}

func (h *BorrowHandler) Route(app *gin.Engine) {
	grp := app.Group(server.RootBorrow)
	grp.GET("/:id", h.getBorrowByID)
	grp.GET("/", h.getList)
	grp.POST("/", h.hr.AuthAccess(), h.create)
	grp.DELETE("/:id", h.hr.AuthAccess(), h.delete)
	grp.PUT("/:id", h.hr.AuthAccess(), h.update)
}

func (h *BorrowHandler) create(c *gin.Context) {
	var req dto.BorrowBookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.hr.BindingError(err)
		return
	}
	err := h.service.CreateBorrow(&req)
	if err != nil {
		h.hr.ErrorInternalServer(c, err)
		return
	}
	c.JSON(http.StatusAccepted, dto.SuccessResponse[any]{
		Success: true,
		Message: "berhasil meminjam buku",
	})
}

func (h *BorrowHandler) getList(c *gin.Context) {
	var req dto.Filter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(h.hr.BindingError(err))
		return
	}
	data, err := h.service.GetList(&req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrDataNotFound):
			c.JSON(http.StatusNotFound, h.hr.ErrorResponse(err.Error()))
		default:
			h.hr.ErrorInternalServer(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[[]dto.BorrowBookRes]{
		Success: true,
		Message: "Daftar Penerbit",
		Data:    data,
	})
}

func (h *BorrowHandler) getBorrowByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	borrow, err := h.service.GetBorrowByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, h.hr.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[dto.BorrowBookRes]{
		Success: true,
		Message: "borrow details",
		Data:    borrow,
	})
}

func (h *BorrowHandler) update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	var req dto.UpdateBorrow
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid data"})
		return
	}
	err = h.service.UpdateBorrow(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update author"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Author updated successfully"})
}

func (h *BorrowHandler) delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	err = h.service.DeleteBorrow(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete borrow"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "borrow deleted successfully"})
}
