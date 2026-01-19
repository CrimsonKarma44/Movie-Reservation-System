package models

import "os"

type Env struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string

	JWTAccessSecret  []byte
	JWTRefreshSecret []byte

	AdminUsername string
	AdminEmail    string
	AdminPassword string
}

func NewEnv() *Env {
	return &Env{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),

		JWTAccessSecret:  []byte(os.Getenv("JWT_ACCESS_SECRET")),
		JWTRefreshSecret: []byte(os.Getenv("JWT_REFRESH_SECRET")),

		AdminUsername: os.Getenv("ADMIN"),
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}
}
