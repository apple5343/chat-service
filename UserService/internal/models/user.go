package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password_hash"`
	Grade     string    `json:"grade" db:"grade"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type UserClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}
