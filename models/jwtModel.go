package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	ID       uint `json:"id"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}
