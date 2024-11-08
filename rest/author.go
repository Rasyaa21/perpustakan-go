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

type AuthorHandler struct {
	hr      *server.Handler
	service *service.AuthorService
}

func NewAuthorHandler(handler *server.Handler, authorService *service.AuthorService) *AuthorHandler {
	return &AuthorHandler{hr: handler, service: authorService}
}

func (h *AuthorHandler) Route(app *gin.Engine) {
	grp := app.Group(server.RootAuthor)
	grp.GET("", h.getList)
	grp.GET("/:id", h.getByID)
	grp.POST("/", h.hr.AuthAccess(), h.create)
	grp.DELETE("/:id", h.hr.AuthAccess(), h.delete)
	grp.PUT("/:id", h.hr.AuthAccess(), h.update)
}

func (h *AuthorHandler) create(c *gin.Context) {
	var req dto.AuthorCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		h.hr.BindingError(err)
		return
	}
	err := h.service.CreateAuthor(&req)
	if err != nil {
		h.hr.ErrorInternalServer(c, err)
		return
	}
	c.JSON(http.StatusAccepted, dto.SuccessResponse[any]{
		Success: true,
		Message: "data berhasil disimpan",
	})
}

func (h *AuthorHandler) getList(c *gin.Context) {
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

	c.JSON(http.StatusOK, dto.SuccessResponse[[]dto.AuthorDetailRes]{
		Success: true,
		Message: "Daftar Penerbit",
		Data:    data,
	})
}

func (h *AuthorHandler) getByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	author, err := h.service.GetAuthorByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, h.hr.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[dto.AuthorDetailRes]{
		Success: true,
		Message: "Author details",
		Data:    author,
	})
}

func (h *AuthorHandler) update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	var req dto.AuthorUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid data"})
		return
	}
	err = h.service.UpdateAuthor(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update author"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Author updated successfully"})
}

func (h *AuthorHandler) delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = h.service.DeleteAuthor(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete author"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Author deleted successfully"})
}
