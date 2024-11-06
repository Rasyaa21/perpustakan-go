package dao

import (
	"base-gin/util"
	"time"
)

type Account struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"size:16;not null;uniqueIndex:user_pass;"`
	Password  string `gorm:"size:255;not null;uniqueIndex:user_pass;"`
}

func NewUser(uname, paswd, secret string) (Account, error) {
	account := Account{
		Username: uname,
	}

	if err := account.SetPassword(paswd, secret); err != nil {
		return Account{}, err
	}

	return account, nil
}

func (t *Account) VerifyPassword(plainPaswd string) bool {
	return util.VerifyPasswordHash(t.Password, plainPaswd)
}

func (t *Account) SetPassword(passsword, secretKey string) error {
	passwordHashed, err := util.PasswordHash(passsword)
	if err != nil {
		return err
	}

	t.Password = passwordHashed

	return nil
}
