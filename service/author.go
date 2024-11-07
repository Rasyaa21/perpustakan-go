package service

import (
	"base-gin/domain"
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/repository"
)

type AuthorService struct {
	repo *repository.AuthorRepository
}

func NewAuthorService(authorRepo *repository.AuthorRepository) *AuthorService {
	return &AuthorService{repo: authorRepo}
}

func (s *AuthorService) CreateAuthor(params *dto.AuthorCreate) error {
	newItem := params.ToEntity()
	return s.repo.Create(&newItem)
}

func (s *AuthorService) GetAuthorByID(id uint) (dto.AuthorDetailRes, error) {
	var res dto.AuthorDetailRes
	item, err := s.repo.GetAuthorByID(id)
	if err != nil {
		return res, err
	}
	if item == nil {
		return res, exception.ErrDataNotFound
	}
	res = dto.AuthorDetailRes{
		ID:       int(item.ID),
		Fullname: item.Fullname,
	}
	return res, nil
}

func (s *AuthorService) GetList(params *dto.Filter) ([]dto.AuthorDetailRes, error) {
	var res []dto.AuthorDetailRes
	items, err := s.repo.GetList(params)
	if err != nil {
		return nil, err
	}
	if len(items) < 1 {
		return nil, exception.ErrDataNotFound
	}
	for _, item := range items {
		var t dto.AuthorDetailRes
		t.FromEntityRes(&item)
		res = append(res, t)
	}
	return res, nil
}

func (s *AuthorService) DeleteAuthor(id uint) error {
	if id <= 0 {
		return exception.ErrDataNotFound
	}
	return s.repo.Delete(id)
}

func (s *AuthorService) UpdateAuthor(id uint, req *dto.AuthorUpdateReq) error {
	author, err := s.repo.GetAuthorByID(id)
	if err != nil {
		return err
	}

	author.Fullname = req.Fullname
	author.Gender = (*domain.TypeGender)(&req.Gender)

	err = s.repo.Update(author)
	if err != nil {
		return err
	}

	return nil
}
