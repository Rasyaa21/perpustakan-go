package dto

import (
	"base-gin/domain"
	"base-gin/domain/dao"
	"time"

	"gorm.io/gorm"
)

type AuthorUpdateReq struct {
	ID           uint      `json:"-"`
	Fullname     string    `json:"fullname" binding:"required,min=4,max=56"`
	Gender       string    `json:"gender" binding:"required,oneof=m f"`
	BirthDateStr string    `json:"birth_date" binding:"required,datetime=2006-01-02"`
	BirthDate    time.Time `json:"-"`
	Age          int       `json:"age"`
}

type AuthorCreate struct {
	gorm.Model
	Fullname     string    `json:"fullname" binding:"required,min=4,max=56"`
	Gender       string    `json:"gender" binding:"required,oneof=m f"`
	BirthDateStr string    `json:"birth_date" binding:"required,datetime=2006-01-02"`
	BirthDate    time.Time `json:"-"`
	Age          int       `json:"age"`
}

type AuthorDetailRes struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

func (o *AuthorCreate) FromEntity(item *dao.Author) {
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
	o.Age = o.ConvertAge()
}

func (o *AuthorCreate) UpdateGender() domain.TypeGender {
	if o.Gender == "f" {
		return domain.GenderFemale
	}
	return domain.GenderMale
}

func (o *AuthorCreate) ConvertAge() int {
	now := time.Now()
	age := now.Year() - o.BirthDate.Year()

	if now.YearDay() < o.BirthDate.YearDay() {
		age--
	}
	return age
}

func (o *AuthorCreate) ToEntity() dao.Author {
	return dao.Author{
		Fullname:  o.Fullname,
		Gender:    (*domain.TypeGender)(&o.Gender),
		BirthDate: &o.BirthDate,
	}
}

func (o *AuthorDetailRes) FromEntityRes(item *dao.Author) {
	o.Fullname = item.Fullname
	o.Gender = string(*item.Gender)
}
