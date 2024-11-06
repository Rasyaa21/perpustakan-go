package integration_test

import (
	"base-gin/domain/dao"
	"base-gin/domain/dto"
	"base-gin/server"
	"base-gin/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublisher_Create_Success(t *testing.T) {
	params := dto.PublisherCreateReq{
		Name: util.RandomStringAlpha(6),
		City: util.RandomStringAlpha(8),
	}

	w := doTest(
		"POST",
		server.RootPublisher,
		params,
		createAuthAccessToken(dummyAdmin.Account.Username),
	)
	assert.Equal(t, 201, w.Code)
}

func TestPublisher_Update_Success(t *testing.T) {
	o := dao.Publisher{
		Name: util.RandomStringAlpha(6),
		City: util.RandomStringAlpha(8),
	}
	_ = publisherRepo.Create(&o)

	params := dto.PublisherUpdateReq{
		Name: util.RandomStringAlpha(7),
		City: util.RandomStringAlpha(10),
	}

	w := doTest(
		"PUT",
		fmt.Sprintf("%s/%d", server.RootPublisher, o.ID),
		params,
		createAuthAccessToken(dummyAdmin.Account.Username),
	)
	assert.Equal(t, 200, w.Code)

	item, _ := publisherRepo.GetByID(o.ID)
	assert.Equal(t, params.Name, item.Name)
	assert.Equal(t, params.City, item.City)
	assert.Equal(t, false, item.DeletedAt.Valid)
}

func TestPublisher_Delete_Success(t *testing.T) {
	o := dao.Publisher{
		Name: util.RandomStringAlpha(6),
		City: util.RandomStringAlpha(8),
	}
	_ = publisherRepo.Create(&o)

	w := doTest(
		"DELETE",
		fmt.Sprintf("%s/%d", server.RootPublisher, o.ID),
		nil,
		createAuthAccessToken(dummyAdmin.Account.Username),
	)
	assert.Equal(t, 200, w.Code)

	item, _ := publisherRepo.GetByID(o.ID)
	assert.Nil(t, item)
}

func TestPublisher_GetList_Success(t *testing.T) {
	o1 := dao.Publisher{
		Name: util.RandomStringAlpha(6),
		City: util.RandomStringAlpha(8),
	}
	_ = publisherRepo.Create(&o1)

	o2 := dao.Publisher{
		Name: util.RandomStringAlpha(6),
		City: util.RandomStringAlpha(8),
	}
	_ = publisherRepo.Create(&o2)

	w := doTest(
		"GET",
		server.RootPublisher,
		nil,
		"",
	)
	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	assert.Contains(t, body, o1.Name)
	assert.Contains(t, body, o2.Name)

	w = doTest(
		"GET",
		server.RootPublisher+"?q="+o1.Name,
		nil,
		"",
	)
	assert.Equal(t, 200, w.Code)

	body = w.Body.String()
	assert.Contains(t, body, o1.Name)
	assert.NotContains(t, body, o2.Name)
}

func TestPublisher_GetDetail_Success(t *testing.T) {
	o := dao.Publisher{
		Name: util.RandomStringAlpha(6),
		City: util.RandomStringAlpha(8),
	}
	_ = publisherRepo.Create(&o)

	w := doTest(
		"GET",
		fmt.Sprintf("%s/%d", server.RootPublisher, o.ID),
		nil,
		"",
	)
	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	assert.Contains(t, body, o.Name)
}
