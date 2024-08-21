package utils

import (
	"github.com/golang-jwt/jwt"
	"os"
)

type UserClaims struct {
	UserId uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func NewAcessToken(claims *UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func NewRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseAcessToken(accessToken string) (*UserClaims, error) {
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

func ParseRefreshToken(refreshToken string) (*jwt.StandardClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*jwt.StandardClaims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
