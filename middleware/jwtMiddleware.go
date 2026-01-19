package middleware

import (
	"context"
	"encoding/json"
	"movie-reservation-system/models"
	"movie-reservation-system/services"
	"movie-reservation-system/utils"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	JwtSecretKeyAccess  []byte
	JwtSecretKeyRefresh []byte

	RefreshStore map[uint]string
	AccessStore  map[uint]string
}

func NewAuthMiddleware(env *models.Env, refreshStore map[uint]string, accessStore map[uint]string) *AuthMiddleware {
	return &AuthMiddleware{
		JwtSecretKeyAccess:  env.JWTAccessSecret,
		JwtSecretKeyRefresh: env.JWTRefreshSecret,
		RefreshStore:        refreshStore,
		AccessStore:         accessStore,
	}
}

func (a *AuthMiddleware) RenewTokenMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "no refresh token", http.StatusUnauthorized)
			return
		}

		refreshToken := cookie.Value
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
			return a.JwtSecretKeyRefresh, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		// Check refresh store
		if stored, ok := a.RefreshStore[claims.ID]; !ok || stored != refreshToken {
			http.Error(w, "refresh token revoked", http.StatusUnauthorized)
			return
		}

		// Generate new access token
		access, newRefresh, err := services.AuthService{}.GenerateToken(claims.ID, claims.IsAdmin, a.JwtSecretKeyAccess, a.JwtSecretKeyRefresh, a.RefreshStore)
		if err != nil {
			http.Error(w, "could not generate new token", http.StatusInternalServerError)
			return
		}

		// Rotate refresh token
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    newRefresh,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
		})

		json.NewEncoder(w).Encode(map[string]string{
			"access_token": access,
		})

		next.ServeHTTP(w, r)
	})
}

func (a *AuthMiddleware) ProtectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("access_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			} else {
				http.Error(w, "Bad Request", http.StatusBadRequest)
			}
			return
		}
		
		token := tokenString.Value

		// Validate the existing token
		claim, err := services.AuthService{}.ValidateJWT(token, a.JwtSecretKeyAccess)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserContextKey, claim)
		// fmt.Println(claim.ID, claim.IsAdmin)

		// Call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
