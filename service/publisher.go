package service

import (
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/repository"
)

type PublisherService struct {
	repo *repository.PublisherRepository
}

func NewPublisherService(publisherRepo *repository.PublisherRepository) *PublisherService {
	return &PublisherService{repo: publisherRepo}
}

func (s *PublisherService) Create(params *dto.PublisherCreateReq) error {
	newItem := params.ToEntity()
	return s.repo.Create(&newItem)
}

func (s *PublisherService) GetByID(id uint) (dto.PublisherResp, error) {
	var resp dto.PublisherResp

	item, err := s.repo.GetByID(id)
	if err != nil {
		return resp, err
	}
	if item == nil {
		return resp, exception.ErrDataNotFound
	}

	resp.FromEntity(item)
	resp.City = item.City

	return resp, nil
}

func (s *PublisherService) GetList(params *dto.Filter) ([]dto.PublisherResp, error) {
	var resp []dto.PublisherResp

	items, err := s.repo.GetList(params)
	if err != nil {
		return nil, err
	}
	if len(items) < 1 {
		return nil, exception.ErrDataNotFound
	}

	for _, item := range items {
		var t dto.PublisherResp
		t.FromEntity(&item)

		resp = append(resp, t)
	}

	return resp, nil
}

func (s *PublisherService) Update(params *dto.PublisherUpdateReq) error {
	if params.ID <= 0 {
		return exception.ErrDataNotFound
	}

	return s.repo.Update(params)
}

func (s *PublisherService) Delete(id uint) error {
	if id <= 0 {
		return exception.ErrDataNotFound
	}

	return s.repo.Delete(id)
}
