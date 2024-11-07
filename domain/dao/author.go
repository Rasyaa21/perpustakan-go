package dao

import (
	"base-gin/domain"
	"time"

	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	Fullname  string             `gorm:"size:56;not null;"`
	Gender    *domain.TypeGender `gorm:"type:enum('f','m');not null;"`
	BirthDate *time.Time
	//1 author has many books
	Book []Book `gorm:"foreignKey:AuthorID"`
}
