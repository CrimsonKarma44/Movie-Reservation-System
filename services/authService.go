package services

import (
	"encoding/json"
	"fmt"
	"movie-reservation-system/models"
	"movie-reservation-system/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	fmt.Println("AuthService initialized")
	return &AuthService{
		db: db,
	}
}

func (authService *AuthService) SignUp(user models.User) ([]byte, error) {
	// SECURITY FIX: Validate email format
	if err := utils.ValidateEmail(user.Email); err != nil {
		return nil, fmt.Errorf("invalid email: %v", err)
	}

	// SECURITY FIX: Normalize email to lowercase for consistent storage and comparison
	user.Email = utils.NormalizeEmail(user.Email)

	var existingUser models.User
	if err := authService.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		response, _ := json.Marshal(map[string]any{
			"id":        existingUser.ID,
			"email":     existingUser.Email,
			"createdAt": existingUser.CreatedAt,
			"message":   "User Already exists",
		})
		return response, nil
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	// fmt.Printf("Attempting to create user: %s\n", user.Email)
	if err := authService.db.Create(&user).Error; err != nil {
		fmt.Printf("Database error creating user: %v\n", err)
		return nil, err
	}

	response, _ := json.Marshal(map[string]any{
		"id":        user.ID,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"message":   "User created successfully",
	})

	return response, nil
}

func (authService *AuthService) Login(user *models.User) ([]byte, error) {
	var existingUser models.User
	// if err := authService.db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
	// 	return nil, err
	// }

	// // TODO: prepare custome error
	// if err := utils.ComparePassword(existingUser.Password, user.Password); err != nil {
	// 	return nil, err
	// }

	if err := authService.db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		return nil, fmt.Errorf("invalid credentials") // Generic message
	}

	if err := utils.ComparePassword(existingUser.Password, user.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials") // Same message
	}

	*user = existingUser

	response, _ := json.Marshal(map[string]any{
		"id":        user.ID,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"message":   "User logged in successfully",
	})

	return response, nil
}

// Updated services/authService.go - GenerateToken function
func (authService AuthService) GenerateToken(id uint, is_admin bool, secret_access []byte, secret_refresh []byte, refreshStore *models.SafeTokenStore) (string, string, error) {
	// Access token (short lived)
	accessClaims := models.Claims{
		ID:      id,
		IsAdmin: is_admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("access:%d", id), // Add subject for identification
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secret_access)
	if err != nil {
		return "", "", err
	}

	// Refresh token (long lived)
	refreshClaims := models.Claims{
		ID:      id,
		IsAdmin: is_admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("refresh:%d", id), // Add subject for identification
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secret_refresh)
	if err != nil {
		return "", "", err
	}

	// Store refresh token securely
	refreshStore.Set(id, refreshToken)

	return accessToken, refreshToken, nil
}

// Update the Logout function to properly revoke tokens
func (authService *AuthService) Logout(userID uint, refreshStore *models.SafeTokenStore) error {
	// Remove the refresh token from store to invalidate it
	refreshStore.Delete(userID)
	return nil
}

// func (authService AuthService) GenerateToken(id uint, is_admin bool, secret_access []byte, secret_refresh []byte, refreshStore map[uint]string) (string, string, error) {
// 	// Access token (short lived)
// 	accessClaims := models.Claims{
// 		ID:       id,
// 		IsAdmin: is_admin,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 		},
// 	}
// 	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secret_access)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	// Refresh token (long lived)
// 	refreshClaims := models.Claims{
// 		ID:       id,
// 		IsAdmin: is_admin,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 		},
// 	}
// 	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secret_refresh)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	// Store refresh token
// 	refreshStore[id] = refreshToken

// 	return accessToken, refreshToken, nil
// }

// func (authService *AuthService) Logout(user *models.User) error {
// 	return nil
// }

func (authService AuthService) ValidateJWT(tokenString string, jwtkey_access []byte) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		// SECURITY FIX: Validate signing algorithm to prevent algorithm substitution attacks
		// Only allow HMAC-based signatures (HS256, HS384, HS512)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v, expected HMAC", token.Header["alg"])
		}
		return jwtkey_access, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	// SECURITY FIX: Explicitly check token expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, jwt.ErrTokenExpired
	}

	return claims, nil
}
