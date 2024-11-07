package service

import (
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/repository"
)

type BookService struct {
	repo *repository.BookRepository
}

func NewBookService(bookRepo *repository.BookRepository) *BookService {
	return &BookService{repo: bookRepo}
}

func (s *BookService) CreateBook(params *dto.BookCreateReq) error {
	newItem := params.ToEntity()
	return s.repo.Create(&newItem)
}

func (s *BookService) GetBookByID(id uint) (dto.BookDetailRes, error) {
	var res dto.BookDetailRes
	item, err := s.repo.GetBookByID(id)
	if err != nil {
		return res, err
	}
	if item == nil {
		return res, exception.ErrDataNotFound
	}
	res = dto.BookDetailRes{
		Title:    item.Title,
		Subtitle: item.Subtitle,
	}
	return res, nil
}

func (s *BookService) GetList(params *dto.Filter) ([]dto.BookDetailRes, error) {
	var res []dto.BookDetailRes
	items, err := s.repo.GetList(params)
	if err != nil {
		return nil, err
	}
	if len(items) < 1 {
		return nil, exception.ErrDataNotFound
	}
	for _, item := range items {
		var t dto.BookDetailRes
		t.BookRes(&item)
		res = append(res, t)
	}
	return res, nil
}

func (s *BookService) DeleteBook(id uint) error {
	if id <= 0 {
		return exception.ErrDataNotFound
	}
	return s.repo.Delete(id)
}

func (s *BookService) UpdateBook(id uint, req *dto.UpdateBook) error {
	book, err := s.repo.GetBookByID(id)
	if err != nil {
		return err
	}
	book.Title = req.Title
	book.Subtitle = req.Subtitle
	err = s.repo.Update(book)
	if err != nil {
		return err
	}
	return nil
}
