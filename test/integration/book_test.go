package integration_test

import (
	"base-gin/domain/dto"
	"base-gin/server"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBook_Create_Success(t *testing.T) {
	req := dto.BookCreateReq{
		Title:    "Sample Book",
		Subtitle: "Lorem12",
		AuthorID: 1,
	}

	w := doTest("POST", server.RootBook, req, createAuthAccessToken(dummyAdmin.Account.Username))
	assert.Equal(t, 202, w.Code)
}

func TestBook_GetList_Success(t *testing.T) {
	w := doTest("GET", server.RootBook, nil, createAuthAccessToken(dummyAdmin.Account.Username))
	assert.Equal(t, 200, w.Code)

	resp := w.Body.String()
	assert.Contains(t, resp, "Sample Book")
}

func TestBook_GetByID_Success(t *testing.T) {
	w := doTest("GET", server.RootBook+"/1", nil, createAuthAccessToken(dummyAdmin.Account.Username))
	assert.Equal(t, 200, w.Code)

	resp := w.Body.String()
	assert.Contains(t, resp, "Sample Book")
}

func TestBook_Update_Success(t *testing.T) {
	req := dto.UpdateBook{
		Title: "Updated Sample Book",
	}

	w := doTest("PUT", server.RootBook+"/1", req, createAuthAccessToken(dummyAdmin.Account.Username))
	assert.Equal(t, 200, w.Code)

	resp := w.Body.String()
	assert.Contains(t, resp, "Updated Sample Book")
}

func TestBook_Delete_Success(t *testing.T) {
	w := doTest("DELETE", server.RootBook+"/1", nil, createAuthAccessToken(dummyAdmin.Account.Username))
	assert.Equal(t, 200, w.Code)

	w = doTest("GET", server.RootBook+"/1", nil, createAuthAccessToken(dummyAdmin.Account.Username))
	assert.Equal(t, 404, w.Code)
}

func TestBook_GetByID_Error_InvalidID(t *testing.T) {
	w := doTest("GET", server.RootBook+"/abc", nil, createAuthAccessToken(dummyAdmin.Account.Username))
	assert.Equal(t, 400, w.Code)
}

func TestBook_Create_Error_InvalidData(t *testing.T) {
	req := dto.BookCreateReq{
		Title: "",
	}

	w := doTest("POST", server.RootBook, req, createAuthAccessToken(dummyAdmin.Account.Username))
	assert.Equal(t, 400, w.Code)
}
