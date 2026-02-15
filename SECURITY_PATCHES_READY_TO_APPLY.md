# Ready-to-Apply Security Patches

This document contains exact code that you can copy-paste to fix the security issues. All patches are ready to apply immediately.

---

## PATCH 1: Remove Password Logging from authHandler.go

**File:** `handlers/authHandler.go`  
**Line:** 50  
**Time to fix:** 2 minutes

### Current Code (REMOVE):
```go
fmt.Printf("Email: %s, Password: %s\n", creds.Email, creds.Password)
```

### Action:
Delete this entire line. It serves no purpose and exposes credentials.

---

## PATCH 2: Fix Login Error Messages - Service Layer

**File:** `services/authService.go`  
**Lines:** 57-72  
**Time to fix:** 10 minutes

### Replace this:
```go
func (authService *AuthService) Login(user *models.User) ([]byte, error) {
	var existingUser models.User
	if err := authService.db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		return nil, err
	}

	// TODO: prepare custome error
	if err := utils.ComparePassword(existingUser.Password, user.Password); err != nil {
		return nil, err
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
```

### With this:
```go
func (authService *AuthService) Login(user *models.User) ([]byte, error) {
	var existingUser models.User
	if err := authService.db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := utils.ComparePassword(existingUser.Password, user.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
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
```

**Changes made:**
- Line 58: Changed generic `err` to `fmt.Errorf("invalid credentials")`
- Line 62: Changed generic `err` to same generic message
- Removed "TODO" comment (already fixed)

---

## PATCH 3: Fix Login Error Response - Handler Layer

**File:** `handlers/authHandler.go`  
**Lines:** 67-77  
**Time to fix:** 5 minutes

### Replace this:
```go
		response, err := h.AuthService.Login(&creds)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not login: %s", err), http.StatusInternalServerError)
			return
		}
```

### With this:
```go
		response, err := h.AuthService.Login(&creds)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
```

**Changes made:**
- Removed error details in message (was: `fmt.Sprintf("could not login: %s", err)`)
- Changed status from `StatusInternalServerError (500)` to `StatusUnauthorized (401)`
- Generic message "invalid credentials" used

---

## PATCH 4: Validate JWT Secrets - Environment Model

**File:** `models/envModel.go`  
**Time to fix:** 10 minutes

### Replace entire file with:
```go
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
	// Validate JWT secrets exist
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	
	if accessSecret == "" {
		log.Fatal("FATAL ERROR: JWT_ACCESS_SECRET environment variable is not set")
	}
	if refreshSecret == "" {
		log.Fatal("FATAL ERROR: JWT_REFRESH_SECRET environment variable is not set")
	}
	
	// Validate minimum secret length (32 characters for HS256 security)
	if len(accessSecret) < 32 {
		log.Fatalf("FATAL ERROR: JWT_ACCESS_SECRET must be at least 32 characters long (currently: %d)", len(accessSecret))
	}
	if len(refreshSecret) < 32 {
		log.Fatalf("FATAL ERROR: JWT_REFRESH_SECRET must be at least 32 characters long (currently: %d)", len(refreshSecret))
	}
	
	return &Env{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),

		JWTAccessSecret:  []byte(accessSecret),
		JWTRefreshSecret: []byte(refreshSecret),

		AdminUsername: os.Getenv("ADMIN"),
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}
}
```

**Changes made:**
- Added validation that secrets are set
- Added minimum length validation (32 characters)
- Added clear error messages explaining the requirements
- Uses `log.Fatal` to prevent app from running with weak secrets

---

## PATCH 5: Validate JWT Algorithm - Auth Service

**File:** `services/authService.go`  
**Lines:** 164-171  
**Time to fix:** 5 minutes

### Replace this:
```go
func (authService AuthService) ValidateJWT(tokenString string, jwtkey_access []byte) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return jwtkey_access, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
```

### With this:
```go
func (authService AuthService) ValidateJWT(tokenString string, jwtkey_access []byte) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		// SECURITY: Validate algorithm to prevent algorithm substitution attacks (e.g., alg=none)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtkey_access, nil
	})
	
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	
	// Explicit validation of claims
	if claims == nil {
		return nil, fmt.Errorf("invalid token claims")
	}
	
	// Verify expiration time is set and valid
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}
	
	return claims, nil
}
```

**Changes made:**
- Added algorithm validation check (only HMAC allowed)
- Added claims nil check
- Added explicit expiration time validation
- Added security comments explaining the checks

**Note:** Make sure to add `"time"` to imports at top of file if not already there.

---

## PATCH 6: Clear Both Tokens on Logout

**File:** `handlers/authHandler.go`  
**Lines:** 132-157 (LogoutHandler function)  
**Time to fix:** 5 minutes

### Replace this:
```go
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println("LogoutHandler called")

		// Get user claims to identify the user
		claims, ok := r.Context().Value(utils.UserContextKey).(*models.Claims)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Remove from store using the safe store
		h.RefreshStore.Delete(claims.ID)

		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})

		w.Write([]byte("logged out"))
	}
}
```

### With this:
```go
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println("LogoutHandler called")

		// Get user claims to identify the user
		claims, ok := r.Context().Value(utils.UserContextKey).(*models.Claims)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Remove from store using the safe store
		h.RefreshStore.Delete(claims.ID)

		secure := os.Getenv("ENV") == "production"

		// Clear access token cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		// Clear refresh token cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "logged out successfully",
		})
	}
}
```

**Changes made:**
- Added clearing of access_token (was missing)
- Added secure flag check
- Added SameSite=Strict to both cookies
- Added JSON response instead of plain text
- Both cookies now have same security settings

---

## PATCH 7: Add Password Validation Utility

**File:** `utils/util.go`  
**Time to fix:** 10 minutes

### Add these functions at the end of the file:

```go
import (
	"fmt"
	"strings"
)

// ValidatePassword checks if password meets security requirements:
// - At least 12 characters long
// - Contains at least one uppercase letter
// - Contains at least one lowercase letter
// - Contains at least one digit
// - Contains at least one special character (!@#$%^&*-_=+)
func ValidatePassword(password string) error {
	if len(password) < 12 {
		return fmt.Errorf("password must be at least 12 characters long")
	}
	
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter (A-Z)")
	}
	
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter (a-z)")
	}
	
	hasNumber := strings.ContainsAny(password, "0123456789")
	if !hasNumber {
		return fmt.Errorf("password must contain at least one digit (0-9)")
	}
	
	hasSpecial := strings.ContainsAny(password, "!@#$%^&*-_=+")
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character (!@#$%%^&*-_=+)")
	}
	
	return nil
}
```

**What to add:**
- Import "strings" if not already imported
- Add the `ValidatePassword` function after existing functions
- This validates password strength according to OWASP guidelines

---

## PATCH 8: Use Password Validation in Registration

**File:** `handlers/authHandler.go`  
**Lines:** 31-58 (RegisterHandler function)  
**Time to fix:** 5 minutes

### Replace the validation section:
```go
		// Validate input
		if creds.Email == "" || creds.Password == "" {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		} else {
			fmt.Printf("Email: %s, Password: %s\n", creds.Email, creds.Password)
		}
```

### With this:
```go
		// Validate input
		if creds.Email == "" || creds.Password == "" {
			http.Error(w, "email and password are required", http.StatusBadRequest)
			return
		}
		
		// Validate email format
		if !strings.Contains(creds.Email, "@") || len(creds.Email) < 5 {
			http.Error(w, "please provide a valid email address", http.StatusBadRequest)
			return
		}
		
		// Validate password strength
		if err := utils.ValidatePassword(creds.Password); err != nil {
			http.Error(w, fmt.Sprintf("password requirement: %s", err.Error()), http.StatusBadRequest)
			return
		}
```

**Changes made:**
- Removed password logging line
- Added email format validation
- Added password strength validation
- Added helpful error messages

**Note:** Add `"strings"` to imports if not present.

---

## PATCH 9: Enhanced Token Store with Expiration

**File:** `models/tokenStoreModel.go`  
**Time to fix:** 15 minutes

### Replace entire file with:
```go
package models

import (
	"sync"
	"time"
)

type TokenEntry struct {
	Token     string
	ExpiresAt time.Time
}

type SafeTokenStore struct {
	store map[uint]TokenEntry
	mutex sync.RWMutex
	done  chan struct{}
}

func NewSafeTokenStore() *SafeTokenStore {
	store := &SafeTokenStore{
		store: make(map[uint]TokenEntry),
		done:  make(chan struct{}),
	}
	
	// Start cleanup goroutine to prevent memory leaks
	go store.cleanupExpiredTokens()
	
	return store
}

func (s *SafeTokenStore) Set(userID uint, token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store[userID] = TokenEntry{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
}

func (s *SafeTokenStore) Get(userID uint) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	entry, exists := s.store[userID]
	if !exists {
		return "", false
	}
	
	// Check if token has expired
	if time.Now().After(entry.ExpiresAt) {
		return "", false
	}
	
	return entry.Token, true
}

func (s *SafeTokenStore) Delete(userID uint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.store, userID)
}

func (s *SafeTokenStore) Validate(userID uint, token string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	entry, exists := s.store[userID]
	if !exists {
		return false
	}
	
	// Check if token has expired
	if time.Now().After(entry.ExpiresAt) {
		return false
	}
	
	return entry.Token == token
}

// cleanupExpiredTokens periodically removes expired tokens to prevent memory leaks
func (s *SafeTokenStore) cleanupExpiredTokens() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			s.mutex.Lock()
			now := time.Now()
			for userID, entry := range s.store {
				if now.After(entry.ExpiresAt) {
					delete(s.store, userID)
				}
			}
			s.mutex.Unlock()
		}
	}
}

// Stop gracefully shuts down the cleanup goroutine
func (s *SafeTokenStore) Stop() {
	close(s.done)
}
```

**Changes made:**
- Added `TokenEntry` struct with expiration time
- Added automatic cleanup goroutine
- Tokens now expire after 24 hours
- Prevents memory leaks from token accumulation
- Thread-safe operation maintained

---

## Verification Commands

After applying patches, run these commands to verify:

### 1. Check password logging removed:
```bash
grep -r "fmt.Printf.*Password" ./handlers/
# Should return nothing (empty result = good)
```

### 2. Check error messages are generic:
```bash
grep -n "could not login" ./handlers/
# Should return nothing (message should be "invalid credentials" now)
```

### 3. Check JWT validation updated:
```bash
grep -n "SigningMethodHMAC" ./services/
# Should show algorithm validation present
```

### 4. Check both tokens cleared:
```bash
grep -n "access_token" ./handlers/authHandler.go
# Should show both access_token and refresh_token being cleared
```

### 5. Check password validation present:
```bash
grep -n "ValidatePassword" ./handlers/authHandler.go
# Should show password validation being called
```

---

## Testing After Patches

### Test 1: Password logging removed
```bash
# Start app and register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"TestPass123!@#"}'

# Check console - should NOT see password printed
```

### Test 2: Generic error messages
```bash
# Wrong email
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"wrong@test.com","password":"TestPass123!@#"}'

# Wrong password
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"WrongPass123!@#"}'

# Both should return: "invalid credentials"
```

### Test 3: JWT secrets enforced
```bash
# Delete JWT_ACCESS_SECRET from .env
# Try to start app
# Should fail with clear error message
```

### Test 4: Password validation works
```bash
# Weak password
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test2@test.com","password":"weak"}'
# Should fail with helpful message

# Strong password
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test2@test.com","password":"StrongPass123!@#"}'
# Should succeed
```

---

## Quick Application Guide

1. **First:** Back up your current files:
   ```bash
   git commit -m "Backup before security patches"
   ```

2. **Apply patches in order:**
   - Patch 1: Remove password logging (2 min)
   - Patch 2-3: Fix error messages (15 min)
   - Patch 4: Validate JWT secrets (10 min)
   - Patch 5: Algorithm validation (5 min)
   - Patch 6: Clear both tokens (5 min)
   - Patch 7-8: Password validation (20 min)
   - Patch 9: Token store improvements (15 min)

3. **Test each patch:**
   - Run tests after each patch
   - Verify the fix works as expected
   - Commit: `git commit -m "Apply security patch X"`

4. **Total time:** ~1.5 hours for all 9 patches

---

*These patches address the 6 critical and high-priority security issues. Apply them before production deployment.*

