package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	UserID       uint   `json:"user_id"`
	UserFullName string `json:"user_full_name"`
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

func CreateTokenWithUserID(userID uint, user_full_name string) (string, error) {
	accessToken, err := NewToken(&UserClaims{
		UserID:       userID,
		UserFullName: user_full_name,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	})
	return accessToken, err
}
