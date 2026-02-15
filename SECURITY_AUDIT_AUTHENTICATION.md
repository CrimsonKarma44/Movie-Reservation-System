# Authentication Security Audit Report
## Movie Reservation System

**Date:** February 15, 2026
**Severity Levels:** 🔴 Critical | 🟠 High | 🟡 Medium | 🟢 Low

---

## Executive Summary

Your authentication implementation has some good practices mixed with critical security vulnerabilities. The biggest concerns are:

1. **🔴 CRITICAL**: Credential exposure in error messages and logs
2. **🔴 CRITICAL**: No rate limiting on authentication endpoints
3. **🟠 HIGH**: JWT secret validation issues
4. **🟠 HIGH**: Login response exposes user information without proper validation
5. **🟡 MEDIUM**: Missing security headers and CSRF protection

---

## Detailed Findings

### 1. 🔴 CRITICAL: Credentials Exposed in Error Messages and Logs

**Location:** `handlers/authHandler.go` (Lines 48-50, 67-77)

**Issue:**
```go
// Line 48-50: Password printed to console
fmt.Printf("Email: %s, Password: %s\n", creds.Email, creds.Password)

// Line 67-77: Error messages expose sensitive details
if err != nil {
    http.Error(w, fmt.Sprintf("could not login: %s", err), http.StatusInternalServerError)
}
```

**Risk:**
- Passwords are logged in plaintext
- Error details may be exposed to client/logs
- Sensitive data could be captured by log aggregation systems

**Recommendation:**
```go
// Remove password logging entirely
// Don't expose detailed error messages to clients
if err != nil {
    http.Error(w, "Authentication failed", http.StatusUnauthorized)
    // Log detailed error server-side only (without credentials)
}
```

---

### 2. 🔴 CRITICAL: No Rate Limiting on Auth Endpoints

**Location:** `handlers/authHandler.go`, `server/app.go`

**Issue:**
- No rate limiting on `/register`, `/login`, or `/logout`
- Attackers can perform brute force attacks, credential stuffing, or DoS

**Risk:**
- Brute force password attacks
- Credential stuffing attacks
- Registration spam
- Denial of Service

**Recommendation:** Implement rate limiting using a package like `github.com/go-chi/chi/middleware` or `github.com/ulule/limiter`:

```go
// Add to go.mod
require github.com/ulule/limiter/v3 v3.11.2

// In server setup
import "github.com/ulule/limiter/v3/drivers/store/memory"

limiter := limiter.New(
    memory.NewStore(),
    limiter.Rate{Limit: 5, Period: 1 * time.Minute},
)

// Apply to sensitive endpoints
mux.Use(limiter.Handler())
```

---

### 3. 🔴 CRITICAL: Login Function Returns User Data Without Verification

**Location:** `services/authService.go` (Lines 56-72)

**Issue:**
```go
func (authService *AuthService) Login(user *models.User) ([]byte, error) {
    var existingUser models.User
    if err := authService.db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
        return nil, err  // Returns error on user not found
    }

    // After password comparison fails, user data still used
    if err := utils.ComparePassword(existingUser.Password, user.Password); err != nil {
        return nil, err
    }

    *user = existingUser  // User ID assigned without returning to caller
    // ...
}
```

**Risk:**
1. **Timing Attack**: Different error messages for "user not found" vs "wrong password" allows attackers to enumerate valid email addresses
2. **User ID Exposure**: The function modifies the pointer without proper validation
3. **Information Disclosure**: Error messages reveal whether account exists

**Recommendation:**
```go
func (authService *AuthService) Login(user *models.User) ([]byte, error) {
    var existingUser models.User

    // Use same error message for both cases
    if err := authService.db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
        return nil, fmt.Errorf("invalid credentials")  // Generic message
    }

    if err := utils.ComparePassword(existingUser.Password, user.Password); err != nil {
        return nil, fmt.Errorf("invalid credentials")  // Same message
    }

    *user = existingUser
    return response, nil
}
```

---

### 4. 🟠 HIGH: JWT Secret Validation Not Enforced

**Location:** `models/envModel.go` (Lines 19-27)

**Issue:**
```go
func NewEnv() *Env {
    return &Env{
        // ...
        JWTAccessSecret:  []byte(os.Getenv("JWT_ACCESS_SECRET")),
        JWTRefreshSecret: []byte(os.Getenv("JWT_REFRESH_SECRET")),
        // ...
    }
}
```

**Problems:**
1. No validation that secrets are set
2. No minimum length requirement
3. Empty secrets will silently fail
4. If `JWT_ACCESS_SECRET` is not set, `[]byte("")` is used - very weak!

**Risk:**
- Attackers can forge tokens with empty or weak secrets
- If secrets aren't configured, system operates with zero security

**Recommendation:**
```go
func NewEnv() *Env {
    accessSecret := os.Getenv("JWT_ACCESS_SECRET")
    refreshSecret := os.Getenv("JWT_REFRESH_SECRET")

    // Validate secrets
    if accessSecret == "" || refreshSecret == "" {
        log.Fatal("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must be set")
    }
    if len(accessSecret) < 32 || len(refreshSecret) < 32 {
        log.Fatal("JWT secrets must be at least 32 characters long")
    }

    return &Env{
        // ...
        JWTAccessSecret:  []byte(accessSecret),
        JWTRefreshSecret: []byte(refreshSecret),
        // ...
    }
}
```

---

### 5. 🟠 HIGH: Weak Token Validation in JWT Parsing

**Location:** `middleware/jwtMiddleware.go` (Lines 79-110)

**Issue:**
```go
func (a *AuthMiddleware) ProtectMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString, err := r.Cookie("access_token")
        if err != nil {
            // Missing validation of cookie attributes
            http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
            return
        }

        // No validation of token expiration
        claim, err := services.AuthService{}.ValidateJWT(token, a.JwtSecretKeyAccess)
        // ...
    })
}
```

**Problems:**
1. No explicit token expiration check in `ValidateJWT`
2. JWT library may not enforce expiration by default
3. No validation of token subject or other claims
4. No check for modified/tampered tokens

**Risk:**
- Expired tokens may still be accepted
- Tokens could be modified if algorithm validation is weak

**Recommendation:**
```go
func (authService AuthService) ValidateJWT(tokenString string, jwtkey_access []byte) (*models.Claims, error) {
    claims := &models.Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
        // Validate algorithm to prevent algorithm substitution
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtkey_access, nil
    })

    if err != nil || !token.Valid {
        return nil, jwt.ErrSignatureInvalid
    }

    // Explicitly check expiration
    if time.Now().Unix() > claims.ExpiresAt.Unix() {
        return nil, jwt.ErrTokenExpired
    }

    return claims, nil
}
```

---

### 6. 🟠 HIGH: Token Storage Security Issue

**Location:** `models/tokenStoreModel.go` (Lines 1-43)

**Issue:**
The refresh token store is an in-memory map that:
1. **Doesn't persist** - tokens lost on application restart, users forced to re-login
2. **Doesn't expire** - tokens stored indefinitely until user logs out
3. **No cleanup** - memory leaks if users never logout
4. **Single instance** - doesn't work in distributed systems

```go
type SafeTokenStore struct {
    store map[uint]string
    mutex sync.RWMutex
}
```

**Risk:**
- Token reuse for 24 hours even after logout (if not manually revoked)
- Memory exhaustion over time
- Not suitable for production multi-instance deployments

**Recommendation:** Use Redis or database-backed token store:
```go
// Use Redis for distributed token storage
func (s *RedisTokenStore) Get(userID uint) (string, bool) {
    val, err := s.client.Get(ctx, fmt.Sprintf("refresh_token:%d", userID)).Result()
    if err == redis.Nil {
        return "", false
    }
    return val, true
}
```

---

### 7. 🟡 MEDIUM: Missing Algorithm Validation in JWT Parser

**Location:** `services/authService.go` (Lines 164-171)

**Issue:**
```go
func (authService AuthService) ValidateJWT(tokenString string, jwtkey_access []byte) (*models.Claims, error) {
    claims := &models.Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
        return jwtkey_access, nil  // No algorithm check!
    })
    // ...
}
```

**Problem:** Doesn't validate the signing algorithm - vulnerable to "alg=none" attacks or algorithm substitution

**Risk:**
- Attacker could claim `alg=none` and bypass signature verification
- Attacker could switch from HS256 to RS256 if not careful

**Recommendation:**
```go
func (token *jwt.Token) (any, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return jwtkey_access, nil
}
```

---

### 8. 🟡 MEDIUM: Missing CORS Configuration

**Location:** `server/app.go`

**Issue:** No CORS (Cross-Origin Resource Sharing) headers configured

**Risk:**
- Cross-site request forgery (CSRF) possible
- Tokens could be stolen via cross-origin requests
- No protection against unauthorized domain access

**Recommendation:**
```go
func (s *Server) Run() {
    // Add CORS middleware
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGINS")},  // e.g., "https://example.com"
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders:   []string{"Content-Type"},
        ExposedHeaders:   []string{},  // Don't expose tokens in headers
        AllowCredentials: true,
        MaxAge:           300,
    })

    mux.Use(c.Handler)
}
```

---

### 9. 🟡 MEDIUM: Insecure Cookie Configuration in Some Cases

**Location:** `handlers/authHandler.go` (Lines 107-129)

**Current Code (Good):**
```go
secure := os.Getenv("ENV") == "production"
http.SetCookie(w, &http.Cookie{
    Name:     "access_token",
    Value:    access,
    HttpOnly: true,
    Secure:   secure,  // Good!
    SameSite: http.SameSiteStrictMode,  // Good!
    Path:     "/",
    MaxAge:   900,
})
```

**Issue:** For non-production environments, `Secure: false` allows tokens over HTTP

**Risk:**
- During development, tokens could be intercepted on unencrypted connections
- Developer mistake could deploy with this setting

**Recommendation:**
```go
secure := os.Getenv("ENV") == "production"
if os.Getenv("ENFORCE_HTTPS") == "true" {
    secure = true
}
```

---

### 10. 🟡 MEDIUM: No Password Complexity Requirements

**Location:** `handlers/authHandler.go` (Lines 42-49)

**Issue:**
```go
if creds.Email == "" || creds.Password == "" {
    http.Error(w, "invalid request", http.StatusBadRequest)
    return
}
// That's it - no complexity check!
```

**Risk:**
- Users can set weak passwords like "123456"
- No minimum length requirement
- No character diversity requirements

**Recommendation:**
```go
func ValidatePassword(password string) error {
    if len(password) < 12 {
        return fmt.Errorf("password must be at least 12 characters")
    }
    if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
        return fmt.Errorf("password must contain uppercase letters")
    }
    if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
        return fmt.Errorf("password must contain lowercase letters")
    }
    if !strings.ContainsAny(password, "0123456789") {
        return fmt.Errorf("password must contain numbers")
    }
    if !strings.ContainsAny(password, "!@#$%^&*") {
        return fmt.Errorf("password must contain special characters")
    }
    return nil
}
```

---

### 11. 🟡 MEDIUM: No Access Token in Logout Handler

**Location:** `handlers/authHandler.go` (Lines 132-157)

**Issue:**
```go
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
    // Only clears refresh_token cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    "",
        // ...
    })
    // Missing: clearing access_token cookie!
}
```

**Risk:**
- Access token still valid in browser after logout
- Could be reused if somehow sent

**Recommendation:**
```go
http.SetCookie(w, &http.Cookie{
    Name:     "access_token",
    Value:    "",
    Expires:  time.Now().Add(-1 * time.Hour),
    HttpOnly: true,
    Path:     "/",
})

http.SetCookie(w, &http.Cookie{
    Name:     "refresh_token",
    Value:    "",
    Expires:  time.Now().Add(-1 * time.Hour),
    HttpOnly: true,
    Path:     "/",
})
```

---

### 12. 🟡 MEDIUM: No Email Verification for Registration

**Location:** `services/authService.go` (Lines 24-51)

**Issue:**
```go
func (authService *AuthService) SignUp(user models.User) ([]byte, error) {
    // No email verification
    // No email validation
    // No check if user provided valid email format
}
```

**Risk:**
- Anyone can register with fake email addresses
- Account takeover possible if emails not verified
- Spam accounts can be created

**Recommendation:**
```go
import "net/mail"

// Validate email format
_, err := mail.ParseAddress(user.Email)
if err != nil {
    return nil, fmt.Errorf("invalid email format")
}

// Send verification email (implement later)
// verificationCode := GenerateRandomCode()
// SendVerificationEmail(user.Email, verificationCode)
```

---

## Summary of Issues by Severity

| Severity | Count | Issues |
|----------|-------|--------|
| 🔴 CRITICAL | 3 | Credential exposure, No rate limiting, Login enumeration |
| 🟠 HIGH | 3 | JWT secret validation, Token validation, Token storage |
| 🟡 MEDIUM | 6 | Algorithm validation, CORS, Cookie config, Password strength, Logout, Email verification |

---

## Quick Fix Priority List

### Immediate (This Week)
1. ✅ Remove password logging from `authHandler.go` line 50
2. ✅ Fix login error messages to be generic in `authService.go`
3. ✅ Add JWT secret validation in `envModel.go`
4. ✅ Add algorithm validation in JWT parser

### Short Term (This Month)
5. ✅ Implement rate limiting on auth endpoints
6. ✅ Add password complexity validation
7. ✅ Clear both tokens on logout
8. ✅ Add CORS configuration

### Medium Term (Next Quarter)
9. ✅ Replace in-memory token store with Redis
10. ✅ Implement email verification for registration
11. ✅ Add comprehensive security logging
12. ✅ Implement account lockout after failed login attempts

---

## Positive Security Practices Found ✓

1. **Password Hashing**: Using `bcrypt` with default cost is good
2. **Token Separation**: Using different secrets for access/refresh tokens
3. **Short-lived Access Tokens**: 15-minute expiration is reasonable
4. **HttpOnly Cookies**: Prevents XSS token theft
5. **SameSite Strict**: CSRF protection on cookies
6. **Token Rotation**: Refresh token rotation implemented
7. **Refresh Token Reuse Detection**: Detecting token replay attacks

---

## Recommended Resources

- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [Rate Limiting Solutions for Go](https://github.com/ulule/limiter)
- [Go JWT Library Security](https://github.com/golang-jwt/jwt)

---

## Next Steps

1. Review each critical issue and implement fixes
2. Add automated security tests
3. Implement comprehensive audit logging
4. Schedule regular security audits
5. Consider penetration testing before production
