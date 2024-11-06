package dto

import (
	"base-gin/domain"
	"base-gin/domain/dao"
	"time"
)

type AccountLoginReq struct {
	Username string `json:"uname" binding:"required,max=16"`
	Password string `json:"paswd" binding:"required,min=8,max=255"`
}

type AccountLoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccountProfileResp struct {
	Fullname string `json:"fullname"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

func (o *AccountProfileResp) FromPerson(person *dao.Person) {
	var gender string
	if person.Gender == nil {
		gender = "-"
	} else if *person.Gender == domain.GenderFemale {
		gender = "wanita"
	} else {
		gender = "pria"
	}

	var age float64
	if person.BirthDate != nil {
		age = time.Since(*person.BirthDate).Hours() / (24 * 365)
	}

	o.Fullname = person.Fullname
	o.Gender = gender
	o.Age = int(age)
}
