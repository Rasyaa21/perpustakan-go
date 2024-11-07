package dto

import (
	"base-gin/domain/dao"

	"gorm.io/gorm"
)

type BookCreateReq struct {
	gorm.Model
	Title       string `json:"title" binding:"required,max=56"`
	Subtitle    string `json:"subtitle" binding:"required,max=64"`
	AuthorID    uint   `json:"author_id" binding:"required"`
	PublisherID uint   `json:"publisher_id" binding:"required"`
}

type BookDetailRes struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	AuthorID    uint   `json:"author_id"`
	PublisherID uint   `json:"publisher_id"`
}

type UpdateBook struct {
	Title    string `json:"title" binding:"required,max=56"`
	Subtitle string `json:"subtitle" binding:"required,max=64"`
}

func (o *BookCreateReq) FromEntity(item *dao.Book) {
	o.Title = item.Title
	o.Subtitle = item.Subtitle
}

func (o *BookDetailRes) BookRes(item *dao.Book) {
	o.Title = item.Title
	o.Subtitle = item.Subtitle
}

func (o *BookCreateReq) ToEntity() dao.Book {
	var item dao.Book
	o.Title = item.Title
	o.Subtitle = item.Subtitle
	return item
}

func (o *UpdateBook) UpdateBook(item *dao.Book, id uint) {
	o.Title = item.Title
	o.Subtitle = item.Subtitle
}
