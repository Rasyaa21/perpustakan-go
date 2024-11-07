package dao

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title       string    `gorm:"size:56;not null;"`
	Subtitle    string    `gorm:"size:64;"`
	AuthorID    uint      `gorm:"not null;"`
	Author      Author    `gorm:"foreignKey:AuthorID;references:ID"`
	PublisherID uint      `gorm:"not null"`
	Publisher   Publisher `gorm:"foreignKey:PublisherID;references:ID"`
}
