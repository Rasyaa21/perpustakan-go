package dao

import (
	"time"

	"gorm.io/gorm"
)

type Borrow struct {
	gorm.Model
	BorrowDate  *time.Time `json:"borrow_date"`
	ReturnDate  *time.Time `json:"return_date"`
	BookId      uint       `json:"book_id"`
	PublisherID uint       `json:"publisher_id"`
	Book        Book       `gorm:"foreignKey:BookId;references:ID"`
	Publisher   Publisher  `gorm:"foreignKey:PublisherID;references:ID"`
}
