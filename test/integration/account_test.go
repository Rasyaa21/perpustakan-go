package integration_test

import (
	"base-gin/domain/dto"
	"base-gin/server"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount_Login_Success(t *testing.T) {
	req := dto.AccountLoginReq{
		Username: "admin",
		Password: password,
	}

	w := doTest("POST", server.RootAccount+server.PathLogin, req, "")
	assert.Equal(t, 200, w.Code)
}

func TestAccount_GetProfile_Success(t *testing.T) {
	accessToken := createAuthAccessToken(dummyAdmin.Account.Username)

	w := doTest("GET", server.RootAccount, nil, accessToken)
	assert.Equal(t, 200, w.Code)

	resp := w.Body.String()
	assert.Contains(t, resp, dummyAdmin.Fullname)
}

func TestAccount_GetProfile_ErrorAccessToken(t *testing.T) {
	w := doTest("GET", server.RootAccount, nil, "")
	assert.Equal(t, 401, w.Code)

	w = doTest("GET", server.RootAccount, nil, "accessToken")
	assert.Equal(t, 401, w.Code)
}
