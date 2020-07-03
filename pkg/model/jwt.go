package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/phuwn/tools/errors"
	"github.com/phuwn/tools/util"
)

// TokenInfo - data model contains user's auth info
type TokenInfo struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

var (
	secretKey           = util.Getenv("SERVER_SECRET_KEY", "")
	ErrInvalidToken     = errors.New("token expired, please log out and log in again")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrBadToken         = errors.New("bad token")
)

// GenerateJWTToken - create an access_token that represents user's session
func GenerateJWTToken(info *TokenInfo, expiresAt int64) (string, error) {
	info.ExpiresAt = expiresAt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, info)
	encryptedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.Customize(500, "failed to sign on token", err)
	}
	return encryptedToken, nil
}

// VerifyUserSession - validates user's access_token and returns user's id if it's verified
func VerifyUserSession(tokenString string) (string, error) {
	claims := TokenInfo{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if !token.Valid {
		return "", ErrInvalidToken
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", ErrInvalidSignature
		}
		return "", ErrBadToken
	}
	if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		return "", ErrInvalidToken
	}
	return claims.UserID, nil
}
