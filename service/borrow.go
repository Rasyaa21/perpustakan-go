package service

import (
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/repository"
)

type BorrowService struct {
	repo *repository.BorrowRepository
}

func NewBorrowService(borrowRepo *repository.BorrowRepository) *BorrowService {
	return &BorrowService{repo: borrowRepo}
}

func (s *BorrowService) CreateBorrow(params *dto.BorrowBookReq) error {
	newItem := params.ToEntity()
	return s.repo.Create(newItem)
}

func (s *BorrowService) GetBorrowByID(id uint) (dto.BorrowBookRes, error) {
	var res dto.BorrowBookRes
	item, err := s.repo.GetBookByID(id)
	if err != nil {
		return res, err
	}
	if item == nil {
		return res, exception.ErrDataNotFound
	}
	res = dto.BorrowBookRes{
		BorrowDate: item.BorrowDate,
		ReturnDate: item.ReturnDate,
	}
	return res, nil
}

func (s *BorrowService) GetList(params *dto.Filter) ([]dto.BorrowBookRes, error) {
	var res []dto.BorrowBookRes
	items, err := s.repo.GetList(params)
	if err != nil {
		return nil, err
	}
	if len(items) < 1 {
		return nil, exception.ErrDataNotFound
	}
	for _, item := range items {
		var t dto.BorrowBookRes
		t.BookRes(&item)
		res = append(res, t)
	}
	return res, nil
}

func (s *BorrowService) DeleteBorrow(id uint) error {
	if id <= 0 {
		return exception.ErrDataNotFound
	}
	return s.repo.Delete(id)
}

func (s *BorrowService) UpdateBorrow(id uint, req *dto.UpdateBorrow) error {
	borrow, err := s.repo.GetBookByID(id)
	if err != nil {
		return err
	}
	borrow.ReturnDate = req.ReturnDate
	err = s.repo.Update(borrow)
	if err != nil {
		return err
	}
	return nil
}
