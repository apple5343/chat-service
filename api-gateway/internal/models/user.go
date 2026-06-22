package models

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}
