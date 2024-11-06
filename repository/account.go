package repository

import (
	"base-gin/domain/dao"
	"base-gin/exception"
	"base-gin/storage"
	"errors"

	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(newItem *dao.Account) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	tx := r.db.WithContext(ctx).Create(&newItem)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *AccountRepository) GetByUsername(uname string) (dao.Account, error) {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	var item dao.Account
	tx := r.db.WithContext(ctx).Where(dao.Account{Username: uname}).
		First(&item)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return item, exception.ErrUserNotFound
		}

		return item, tx.Error
	}

	return item, nil
}
