package utils

import (
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

type UserClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func NewToken(claims *UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseToken(accessToken string) (*UserClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*UserClaims)

	if !ok {
		return nil, err
	}

	return claims, nil
}

func CreateTokenWithUserID(userID uint) (string, error) {
	accessToken, err := NewToken(&UserClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	})
	return accessToken, err
}
