package middleware

import (
	"fmt"
	"net/http"
)

// SecurityHeadersMiddleware adds important security headers to all responses
type SecurityHeadersMiddleware struct{}

// NewSecurityHeadersMiddleware creates a new security headers middleware
func NewSecurityHeadersMiddleware() *SecurityHeadersMiddleware {
	return &SecurityHeadersMiddleware{}
}

// Handler returns a middleware function that adds security headers
func (s *SecurityHeadersMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// SECURITY FIX: Prevent clickjacking attacks
			// Only allow the site to be embedded in frames from the same origin
			w.Header().Set("X-Frame-Options", "SAMEORIGIN")

			// SECURITY FIX: Prevent MIME type sniffing
			// Forces browser to respect the Content-Type header
			w.Header().Set("X-Content-Type-Options", "nosniff")

			// SECURITY FIX: Enable XSS protection in older browsers
			// Modern browsers use Content-Security-Policy instead
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			// SECURITY FIX: Strict Transport Security (HTTPS only)
			// Uncomment for production (requires HTTPS)
			// w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

			// SECURITY FIX: Referrer Policy - controls what referrer info is sent
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// SECURITY FIX: Content Security Policy - prevents inline scripts and external resource loading
			w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'")

			// SECURITY FIX: Permissions Policy (formerly Feature Policy) - controls browser features
			w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=(), usb=(), magnetometer=(), gyroscope=(), accelerometer=()")

			// Don't expose server information
			w.Header().Del("Server")

			// SECURITY FIX: Cache control for sensitive pages
			// Prevent caching of authenticated responses
			if isAuthenticatedRequest(r) {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
			}

			next.ServeHTTP(w, r)
		})
	}
}

// isAuthenticatedRequest checks if the request appears to be authenticated
func isAuthenticatedRequest(r *http.Request) bool {
	// Check for presence of authentication token in cookie
	_, err := r.Cookie("access_token")
	if err == nil {
		return true
	}

	// Check for authorization header
	if r.Header.Get("Authorization") != "" {
		return true
	}

	return false
}

// LogSecurityHeaders logs security configuration on startup
func (s *SecurityHeadersMiddleware) LogSecurityHeaders() {
	fmt.Println("✓ Security Headers Configuration:")
	fmt.Println("  - X-Frame-Options: SAMEORIGIN (Clickjacking protection)")
	fmt.Println("  - X-Content-Type-Options: nosniff (MIME sniffing prevention)")
	fmt.Println("  - X-XSS-Protection: enabled (XSS protection)")
	fmt.Println("  - Content-Security-Policy: enabled (Script injection prevention)")
	fmt.Println("  - Referrer-Policy: strict-origin-when-cross-origin")
	fmt.Println("  - Permissions-Policy: restrictive (Feature control)")
}
