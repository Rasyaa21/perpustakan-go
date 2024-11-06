package dto

import (
	"base-gin/domain"
	"base-gin/domain/dao"
	"time"
)

type AuthorUpdateReq struct {
	ID           uint      `json:"-"`
	Fullname     string    `json:"fullname" binding:"required,min=4,max=56"`
	Gender       string    `json:"gender" binding:"required,oneof=m f"`
	BirthDateStr string    `json:"birth_date" binding:"required,datetime=2006-01-02"`
	BirthDate    time.Time `json:"-"`
}

type AuhtorDetailRes struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

func (o *AuthorUpdateReq) FromEntity(item *dao.Author) {
	var gender string
	if item.Gender == nil {
		gender = "-"
	} else if *item.Gender == domain.GenderFemale {
		gender = "Female"
	} else {
		gender = "male"
	}
	o.ID = uint(item.ID)
	o.Gender = gender
	o.Fullname = item.Fullname
	o.BirthDate = *item.BirthDate
}

func (o *AuthorUpdateReq) UpdateGender() domain.TypeGender {
	if o.Gender == "f" {
		return domain.GenderFemale
	}
	return domain.GenderMale
}
