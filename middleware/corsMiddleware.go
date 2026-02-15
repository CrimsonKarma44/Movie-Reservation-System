package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// CORSMiddleware implements CORS (Cross-Origin Resource Sharing) headers
// to prevent cross-site request forgery and control which origins can access the API
type CORSMiddleware struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// NewCORSMiddleware creates a new CORS middleware instance
func NewCORSMiddleware() *CORSMiddleware {
	// Get allowed origins from environment variable
	// Format: "https://example.com,https://app.example.com"
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	if originsEnv == "" {
		// Default to current domain only in production
		if os.Getenv("ENV") == "production" {
			originsEnv = os.Getenv("API_URL")
			if originsEnv == "" {
				originsEnv = "https://localhost:3000" // Fallback
			}
		} else {
			// Allow localhost for development
			originsEnv = "http://localhost:3000,http://localhost:3001,http://localhost:8080"
		}
	}

	origins := strings.Split(originsEnv, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return &CORSMiddleware{
		AllowedOrigins: origins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
	}
}

// Handler returns a middleware function for CORS
func (c *CORSMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get origin from request
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			if c.isOriginAllowed(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Set allowed methods
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.AllowedMethods, ", "))

			// Set allowed headers
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.AllowedHeaders, ", "))

			// Set max age for preflight cache (5 minutes)
			w.Header().Set("Access-Control-Max-Age", "300")

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// isOriginAllowed checks if the origin is in the allowed list
func (c *CORSMiddleware) isOriginAllowed(origin string) bool {
	for _, allowedOrigin := range c.AllowedOrigins {
		if allowedOrigin == origin || allowedOrigin == "*" {
			return true
		}
	}
	return false
}

// LogCORSInfo logs CORS configuration on startup
func (c *CORSMiddleware) LogCORSInfo() {
	fmt.Println("✓ CORS Configuration:")
	fmt.Printf("  Allowed Origins: %v\n", c.AllowedOrigins)
	fmt.Printf("  Allowed Methods: %s\n", strings.Join(c.AllowedMethods, ", "))
	fmt.Printf("  Allowed Headers: %s\n", strings.Join(c.AllowedHeaders, ", "))
}
