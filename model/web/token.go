package web

import "github.com/golang-jwt/jwt"

type AccessClaims struct {
	jwt.StandardClaims
	Email  string `json:"email"`
	RoleID uint   `json:"role_id"`
}

type RefreshClaims struct {
	jwt.StandardClaims
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
