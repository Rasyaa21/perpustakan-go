package dto

import "base-gin/domain/dao"

type PublisherCreateReq struct {
	Name string `json:"name" binding:"required,min=2,max=48"`
	City string `json:"city" binding:"required,max=32"`
}

func (o *PublisherCreateReq) ToEntity() dao.Publisher {
	var item dao.Publisher
	item.Name = o.Name
	item.City = o.City

	return item
}

type PublisherResp struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city,omitempty"`
}

func (o *PublisherResp) FromEntity(item *dao.Publisher) {
	o.ID = int(item.ID)
	o.Name = item.Name
}

type PublisherUpdateReq struct {
	ID   uint   `json:"-"`
	Name string `json:"name" binding:"required,min=2,max=48"`
	City string `json:"city" binding:"required,max=32"`
}
