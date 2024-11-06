package dto

import (
	"base-gin/domain"
	"base-gin/domain/dao"
	"time"
)

type PersonDetailResp struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

func (o *PersonDetailResp) FromEntity(item *dao.Person) {
	var gender string
	if item.Gender == nil {
		gender = "-"
	} else if *item.Gender == domain.GenderFemale {
		gender = "wanita"
	} else {
		gender = "pria"
	}

	var age float64
	if item.BirthDate != nil {
		age = time.Since(*item.BirthDate).Hours() / (24 * 365)
	}

	o.Fullname = item.Fullname
	o.Gender = gender
	o.Age = int(age)
	o.ID = int(item.ID)
}

type PersonUpdateReq struct {
	ID           uint      `json:"-"`
	Fullname     string    `json:"fullname" binding:"required,min=4,max=56"`
	Gender       string    `json:"gender" binding:"required,oneof=m f"`
	BirthDateStr string    `json:"birth_date" binding:"required,datetime=2006-01-02"`
	BirthDate    time.Time `json:"-"`
}

func (o *PersonUpdateReq) GetGender() domain.TypeGender {
	if o.Gender == "f" {
		return domain.GenderFemale
	}

	return domain.GenderMale
}

func (o *PersonUpdateReq) GetBirthDate() (time.Time, error) {
	return time.Parse("2006-01-02", o.BirthDateStr)
}
