package repository

import (
	"base-gin/domain/dao"
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/storage"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type BorrowRepository struct {
	db *gorm.DB
}

func NewBorrowRepository(db *gorm.DB) *BorrowRepository {
	return &BorrowRepository{db: db}
}

func (r *BorrowRepository) Create(newItem *dao.Borrow) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()
	tx := r.db.WithContext(ctx).Create(&newItem)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *BorrowRepository) GetList(params *dto.Filter) ([]dao.Borrow, error) {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	var items []dao.Borrow
	tx := r.db.WithContext(ctx)
	if params.Keyword != "" {
		q := fmt.Sprintln("%%%s%%", params.Keyword)
		tx = tx.Where("Borrow LIKE ?", q)
	}
	if params.Start >= 0 {
		tx.Offset(params.Start)
	}
	if params.Limit > 0 {
		tx.Limit(params.Limit)
	}
	tx = tx.Order("fullname ASC").Find(&items)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, tx.Error
	}
	return items, nil
}

func (r *BorrowRepository) GetBookByID(id uint) (*dao.Borrow, error) {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	var item dao.Borrow
	tx := r.db.WithContext(ctx).First(&item, id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}
	}
	return &item, nil
}

func (r *BorrowRepository) Update(borrow *dao.Borrow) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	tx := r.db.WithContext(ctx).Model(&dao.Borrow{}).Where("id = ?", borrow.ID).Updates(map[string]interface{}{
		"return_date": borrow.ReturnDate,
	})

	return tx.Error
}

func (r *BorrowRepository) Delete(id uint) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()
	var BorrowData dao.Borrow
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&BorrowData).Error; err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&dao.Borrow{}).Error; err != nil {
		return err
	}
	return nil
}
