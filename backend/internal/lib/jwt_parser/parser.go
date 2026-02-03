package jwt_parser

import (
	"fmt"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId int64  `json:"uid"`
	Email  string `json:"email"`
	AppId  int64  `json:"app_id"`
	Exp    int64  `json:"exp"`
	jwt.RegisteredClaims
}

func Parse(tokenStr string, appSecret string) (*Claims, error) {
	claims := Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(appSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while parsing token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return &claims, nil
}
