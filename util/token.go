package util

import (
	"base-gin/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const tokenIssuer = "plus.quranbest.com"

var (
	ErrTokenUnknown               = errors.New("token tidak dikenali")
	ErrTokenVerificationFailed    = errors.New("gagal melakukan verifikasi token")
	ErrTokenInvalid               = errors.New("token tidak valid")
	ErrAuthTokenExpired           = errors.New("token kedaluwarsa")
	ErrAccessTokenFailedToIssue   = errors.New("gagal menerbitkan token access")
	ErrRefreshTokenFailedToIssue  = errors.New("gagal menerbitkan token refresh")
	ErrAccessTokenFailedToVerify  = errors.New("gagal verifikasi token access")
	ErrRefreshTokenFailedToVerify = errors.New("gagal verifikasi token refresh")
)

type AuthAccessClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func CreateAuthAccessToken(cfg config.Config, subject string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthAccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().
				Add(time.Duration(cfg.AuthN.JWTAuthTTL) * time.Second),
			),
			Issuer:   tokenIssuer,
			Audience: jwt.ClaimStrings{"access"},
		},
	})

	signedToken, err := token.SignedString([]byte(cfg.AuthN.JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrAccessTokenFailedToIssue, err)
	}
	return signedToken, nil
}

func CreateAuthRefreshToken(cfg config.Config, subject string) (string, error) {
	refreshClaims := &jwt.RegisteredClaims{
		Subject: subject,
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().
			Add(time.Duration(cfg.AuthN.JWTRefreshTTL) * time.Second),
		),
		Issuer:   tokenIssuer,
		Audience: jwt.ClaimStrings{"refresh"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedRefreshToken, err := token.SignedString([]byte(cfg.AuthN.JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrRefreshTokenFailedToIssue, err.Error())
	}

	return signedRefreshToken, nil
}

func verifyAuthToken(cfg config.Config, authToken string, tokenAud string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %s", ErrTokenUnknown, "signature not match")
		}
		return []byte(cfg.AuthN.JWTSecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrAuthTokenExpired
	}

	accessClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrTokenUnknown, "invalid structure")
	}

	if !accessClaims.VerifyIssuer(tokenIssuer, true) ||
		!accessClaims.VerifyAudience(tokenAud, true) {
		return nil, ErrTokenUnknown
	}

	return accessClaims, nil
}

func VerifyAuthAccessToken(cfg config.Config, token string) (jwt.MapClaims, error) {
	accessClaims, err := verifyAuthToken(cfg, token, "access")
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrAccessTokenFailedToVerify, err.Error())
	}

	return accessClaims, nil
}

func VerifyAuthRefreshToken(cfg config.Config, token string) (jwt.MapClaims, error) {
	accessClaims, err := verifyAuthToken(cfg, token, "refresh")
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRefreshTokenFailedToVerify, err.Error())
	}

	return accessClaims, nil
}
