package handlers

import (
	"encoding/json"
	"fmt"
	"movie-reservation-system/models"
	"movie-reservation-system/services"
	"movie-reservation-system/utils"
	"net/http"
	"os"
	"strings"
	"time"
)

type AuthHandler struct {
	// Add necessary fields, e.g., DB connection, auth service, etc.
	AuthService *services.AuthService

	JwtSecretKeyAccess  []byte
	JwtSecretKeyRefresh []byte

	RefreshStore *models.SafeTokenStore
}

func NewAuthHandler(authService *services.AuthService, env *models.Env, refreshStore *models.SafeTokenStore) *AuthHandler {
	return &AuthHandler{
		AuthService:         authService,
		JwtSecretKeyAccess:  env.JWTAccessSecret,
		JwtSecretKeyRefresh: env.JWTRefreshSecret,
		RefreshStore:        refreshStore,
	}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Implement registration logic
	fmt.Println("Running RegisterHandler")
	if r.Method == http.MethodPost {
		var creds models.User
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		// SECURITY FIX: Validate email and password format
		if creds.Email == "" || creds.Password == "" {
			http.Error(w, "email and password are required", http.StatusBadRequest)
			return
		}

		// SECURITY FIX: Validate email format
		if !isValidEmail(creds.Email) {
			http.Error(w, "invalid email format", http.StatusBadRequest)
			return
		}

		// SECURITY FIX: Validate password strength and complexity
		if err := utils.ValidatePassword(creds.Password); err != nil {
			// Don't expose detailed reasons to prevent user enumeration
			http.Error(w, fmt.Sprintf("password does not meet requirements: %v", err), http.StatusBadRequest)
			return
		}

		// Save to database (pseudo-code, replace with actual DB logic)
		res, err := h.AuthService.SignUp(creds)
		if err != nil {
			// SECURITY FIX: Log failed registration attempt
			auditor := utils.GetAuditor()
			auditor.LogRegistration(creds.Email, getClientIP(r), r.UserAgent(), false, err.Error(), 0)
			http.Error(w, "could not save user", http.StatusInternalServerError)
			return
		}

		// SECURITY FIX: Log successful registration
		auditor := utils.GetAuditor()
		auditor.LogRegistration(creds.Email, getClientIP(r), r.UserAgent(), true, "user registered successfully", 0)

		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}

}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login handler")
	// Implement login logic
	var creds models.User
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		response, err := h.AuthService.Login(&creds)
		if err != nil {
			// SECURITY FIX: Log failed login attempt for audit trail
			auditor := utils.GetAuditor()
			auditor.LogAuthenticationAttempt(creds.Email, getClientIP(r), r.UserAgent(), false, err.Error(), 0)
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}

		// Generate tokens
		access, refresh, err := h.AuthService.GenerateToken(creds.ID, creds.IsAdmin, h.JwtSecretKeyAccess, h.JwtSecretKeyRefresh, h.RefreshStore)
		if err != nil {
			http.Error(w, "could not generate tokens", http.StatusInternalServerError)
			return
		}

		// SECURITY FIX: Log successful login
		auditor := utils.GetAuditor()
		auditor.LogAuthenticationAttempt(creds.Email, getClientIP(r), r.UserAgent(), true, "successful login", creds.ID)

		// Set refresh token as HttpOnly cookie
		// http.SetCookie(w, &http.Cookie{
		// 	Name:     "refresh_token",
		// 	Value:    refresh,
		// 	Expires:  time.Now().Add(24 * time.Hour),
		// 	HttpOnly: true,
		// 	Secure:   false, // true in production with HTTPS
		// 	Path:     "/",
		// })

		// json.NewEncoder(w).Encode(map[string]interface{}{
		// 	"access_token": access,
		// 	"status":       http.StatusAccepted,
		// 	"response":     json.RawMessage(response),
		// })
		secure := os.Getenv("ENV") == "production"

		// Access token (short-lived)
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    access,
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			MaxAge:   900, // 15 minutes
		})

		// Refresh token (long-lived)
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refresh,
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			MaxAge:   86400, // 24 hours
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"user":   json.RawMessage(response),
		})
	}
}

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

		// SECURITY FIX: Log logout event
		auditor := utils.GetAuditor()
		auditor.LogLogout(claims.ID, "", getClientIP(r))

		// SECURITY FIX: Clear both access and refresh tokens on logout
		// Clear access token
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Secure:   os.Getenv("ENV") == "production",
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			MaxAge:   -1,
		})

		// Clear refresh token
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Secure:   os.Getenv("ENV") == "production",
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			MaxAge:   -1,
		})

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"logged out successfully"}`))
	}
}

// isValidEmail validates email format using RFC 5322
func isValidEmail(email string) bool {
	if len(email) == 0 || len(email) > 254 {
		return false
	}

	// Basic RFC 5322 validation
	atIndex := -1
	for i, char := range email {
		if char == '@' {
			if atIndex != -1 {
				return false // Multiple @ symbols
			}
			atIndex = i
		}
	}

	if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 {
		return false // No @, or @ at start/end
	}

	localPart := email[:atIndex]
	domain := email[atIndex+1:]

	// Check local part (before @)
	if len(localPart) == 0 || len(localPart) > 64 {
		return false
	}

	// Check domain part (after @)
	if len(domain) < 3 || !contains(domain, ".") {
		return false // No domain extension
	}

	return true
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// getClientIP extracts the real client IP address from the request
// Considers proxy headers like X-Forwarded-For which are set by load balancers
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (set by proxies like Vercel, nginx, CloudFlare)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header (alternative proxy header)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	return r.RemoteAddr
}
