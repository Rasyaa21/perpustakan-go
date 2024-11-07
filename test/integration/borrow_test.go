package integration_test

import (
	"base-gin/domain/dto"
	"base-gin/server"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBorrow_Create_Success(t *testing.T) {
	borrowDate, _ := time.Parse("2006-01-02", "2024-11-01")
	returnDate, _ := time.Parse("2006-01-02", "2024-11-15")
	req := dto.BorrowBookReq{
		BookId:      1,
		PublisherID: 1,
		BorrowDate:  &borrowDate,
		ReturnDate:  &returnDate,
	}

	w := doTest("POST", server.RootBorrow, req, createAuthAccessToken(dummyAdmin.Account.Username))

	assert.Equal(t, 202, w.Code)
}

func TestBorrow_GetList_Success(t *testing.T) {
	accessToken := createAuthAccessToken(dummyAdmin.Account.Username)
	w := doTest("GET", server.RootBorrow, nil, accessToken)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "borrowed")
}

func TestBorrow_GetByID_Success(t *testing.T) {
	accessToken := createAuthAccessToken(dummyAdmin.Account.Username)
	w := doTest("GET", server.RootBorrow+"/1", nil, accessToken)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "borrow details")
}

func TestBorrow_Update_Success(t *testing.T) {
	returnDate, _ := time.Parse("2006-01-02", "2024-11-15")
	req := dto.UpdateBorrow{
		ReturnDate: &returnDate,
	}
	accessToken := createAuthAccessToken(dummyAdmin.Account.Username)
	w := doTest("PUT", server.RootBook+"/1", req, accessToken)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "borrow updated successfully")
}

func TestBorrow_Delete_Success(t *testing.T) {
	accessToken := createAuthAccessToken(dummyAdmin.Account.Username)
	w := doTest("DELETE", server.RootBook+"/1", nil, accessToken)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "borrow deleted successfully")
}

func TestBorrow_Error_InvalidAccessToken(t *testing.T) {
	w := doTest("GET", server.RootBook, nil, "")
	assert.Equal(t, 401, w.Code)
	w = doTest("GET", server.RootBook, nil, "invalidToken")
	assert.Equal(t, 401, w.Code)
}
