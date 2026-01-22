package middleware

import (
	"movie-reservation-system/models"
	"movie-reservation-system/utils"
	"net/http"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the claims from the token
		claims, ok := r.Context().Value(utils.UserContextKey).(*models.Claims)
		// fmt.Println(claims.ID, claims.IsAdmin)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the user is an admin
		if !claims.IsAdmin {
			http.Error(w, "forbidden: requires admin privileges", http.StatusForbidden)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
