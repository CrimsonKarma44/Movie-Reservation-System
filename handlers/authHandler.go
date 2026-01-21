package handlers

import (
	"encoding/json"
	"fmt"
	"movie-reservation-system/models"
	"movie-reservation-system/services"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	// Add necessary fields, e.g., DB connection, auth service, etc.
	AuthService *services.AuthService

	JwtSecretKeyAccess  []byte
	JwtSecretKeyRefresh []byte

	RefreshStore map[uint]string
}

func NewAuthHandler(authService *services.AuthService, env *models.Env, refreshStore map[uint]string) *AuthHandler {
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
		fmt.Println(creds)

		// Validate input
		if creds.Email == "" || creds.Password == "" {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		} else {
			fmt.Printf("Email: %s, Password: %s\n", creds.Email, creds.Password)
		}

		// Save to database (pseudo-code, replace with actual DB logic)
		res, err := h.AuthService.SignUp(creds)
		if err != nil {
			http.Error(w, "could not save user", http.StatusInternalServerError)
			return
		}

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
			http.Error(w, fmt.Sprintf("could not login: %s", err), http.StatusInternalServerError)
			return
		}

		// Generate tokens
		access, refresh, err := h.AuthService.GenerateToken(creds.ID, creds.IsAdmin, h.JwtSecretKeyAccess, h.JwtSecretKeyRefresh, h.RefreshStore)
		if err != nil {
			http.Error(w, "could not generate tokens", http.StatusInternalServerError)
			return
		}

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
	// Implement logout logic
	if r.Method == http.MethodPost {
		fmt.Println("LogoutHandler called")
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "no refresh token", http.StatusBadRequest)
			return
		}

		refreshToken := cookie.Value
		claims := &models.Claims{}
		jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
			return h.JwtSecretKeyRefresh, nil
		})

		// Remove from store
		delete(h.RefreshStore, claims.ID)

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
