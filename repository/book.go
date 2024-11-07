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

type BookRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(newItem *dao.Book) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()
	tx := r.db.WithContext(ctx).Create(&newItem)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *BookRepository) GetList(params *dto.Filter) ([]dao.Book, error) {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	var items []dao.Book
	tx := r.db.WithContext(ctx)
	if params.Keyword != "" {
		q := fmt.Sprintln("%%%s%%", params.Keyword)
		tx = tx.Where("Book LIKE ?", q)
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

func (r *BookRepository) GetBookByID(id uint) (*dao.Book, error) {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	var item dao.Book
	tx := r.db.WithContext(ctx).First(&item, id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}
	}
	return &item, nil
}

func (r *BookRepository) Update(author *dao.Book) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	tx := r.db.WithContext(ctx).Model(&dao.Book{}).Where("id = ?", author.ID).Updates(map[string]interface{}{
		"Title":  author.Title,
		"gender": author.Subtitle,
	})
	return tx.Error
}

func (r *BookRepository) Delete(id uint) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()
	var book dao.Book
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&book).Error; err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&dao.Book{}).Error; err != nil {
		return err
	}
	return nil
}
