package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("abcdasehfkbdhkasnajb12421ddsswifn001")

type JWTClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
