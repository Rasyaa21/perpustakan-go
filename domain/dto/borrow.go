package dto

import (
	"base-gin/domain/dao"
	"time"

	"gorm.io/gorm"
)

type BorrowBookReq struct {
	gorm.Model
	BorrowDate  *time.Time `json:"borrow_date"`
	ReturnDate  *time.Time `json:"return_date"`
	BookId      uint       `json:"book_id" gorm:"not null;"`
	PublisherID uint       `json:"publisher_id" gorm:"not null"`
}

type BorrowBookRes struct {
	BorrowDate  *time.Time `json:"borrow_date"`
	ReturnDate  *time.Time `json:"return_date"`
	BookId      uint       `json:"book_id" gorm:"not null;"`
	PublisherID uint       `json:"publisher_id" gorm:"not null"`
}

type UpdateBorrow struct {
	ReturnDate *time.Time `json:"return_date"`
}

func (o *BorrowBookReq) FromEntity(item *dao.Borrow) {
	o.BorrowDate = &item.CreatedAt
	o.ReturnDate = item.ReturnDate
}

func (o *BorrowBookReq) ToEntity() *dao.Borrow {
	return &dao.Borrow{
		BorrowDate: o.BorrowDate,
		ReturnDate: o.ReturnDate,
	}
}

func (o *BorrowBookRes) BookRes(item *dao.Borrow) {
	o.BorrowDate = &item.Book.CreatedAt
	o.ReturnDate = nil
}

func (o *UpdateBorrow) UpdateBooks(id uint) {
	o.ReturnDate = &time.Time{}
	*o.ReturnDate = time.Now()
}
