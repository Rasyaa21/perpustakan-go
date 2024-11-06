package dao

import (
	"base-gin/domain"
	"time"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	AccountID *uint              `gorm:"uniqueIndex;"`
	Account   *Account           `gorm:"foreignKey:AccountID;"`
	Fullname  string             `gorm:"size:56;not null;"`
	Gender    *domain.TypeGender `gorm:"type:enum('f','m');"`
	BirthDate *time.Time
}

func (Person) TableName() string {
	return "persons"
}
