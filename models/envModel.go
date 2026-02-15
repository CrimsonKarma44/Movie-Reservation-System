package models

import (
	"fmt"
	"log"
	"os"
)

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
	// Get JWT secrets with fallback to alternative env var names
	jwtAccessSecret := os.Getenv("JWT_SECRET_KEY_ACCESS")
	if jwtAccessSecret == "" {
		jwtAccessSecret = os.Getenv("JWT_ACCESS_SECRET")
	}

	jwtRefreshSecret := os.Getenv("JWT_SECRET_KEY_REFRESH")
	if jwtRefreshSecret == "" {
		jwtRefreshSecret = os.Getenv("JWT_REFRESH_SECRET")
	}

	// SECURITY FIX: Validate JWT secrets
	if jwtAccessSecret == "" || jwtRefreshSecret == "" {
		log.Fatal("SECURITY ERROR: JWT_SECRET_KEY_ACCESS and JWT_SECRET_KEY_REFRESH must be set")
	}

	// SECURITY FIX: Enforce minimum secret length (32 characters recommended)
	const minSecretLength = 32
	if len(jwtAccessSecret) < minSecretLength {
		log.Fatalf("SECURITY ERROR: JWT_SECRET_KEY_ACCESS must be at least %d characters (got %d)", minSecretLength, len(jwtAccessSecret))
	}
	if len(jwtRefreshSecret) < minSecretLength {
		log.Fatalf("SECURITY ERROR: JWT_SECRET_KEY_REFRESH must be at least %d characters (got %d)", minSecretLength, len(jwtRefreshSecret))
	}

	// SECURITY FIX: Ensure secrets are different
	if jwtAccessSecret == jwtRefreshSecret {
		log.Fatal("SECURITY ERROR: JWT_SECRET_KEY_ACCESS and JWT_SECRET_KEY_REFRESH must be different")
	}

	fmt.Println("✓ JWT secrets validated successfully")

	return &Env{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),

		JWTAccessSecret:  []byte(jwtAccessSecret),
		JWTRefreshSecret: []byte(jwtRefreshSecret),

		AdminUsername: os.Getenv("ADMIN"),
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}
}
