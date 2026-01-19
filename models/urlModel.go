package models

import (
	"gorm.io/gorm"
)

type URLModel struct {
	DB  *gorm.DB
	ENV Env

	JWTAccessSecret  string
	JWTRefreshSecret string
}

func NewURLModel(db *gorm.DB, env Env, jwtAccessSecret, jwtRefreshSecret string) *URLModel {
	return &URLModel{
		DB:               db,
		ENV:              env,
		JWTAccessSecret:  jwtAccessSecret,
		JWTRefreshSecret: jwtRefreshSecret,
	}
}

func (u *URLModel) Start() error {
	return nil
}
