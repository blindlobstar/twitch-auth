package tokens

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"

	"github.com/golang-jwt/jwt/v5"
)

type AccessToken struct {
	secret hash.Hash
}

func New(secret string) *AccessToken {
	h := hmac.New(sha256.New, []byte(secret))
	return &AccessToken{
		secret: h,
	}
}

func (token *AccessToken) CreateAccessToken(userId string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
	})

	tokenString, err := at.SignedString([]byte("secret"))
	return tokenString, err
}
