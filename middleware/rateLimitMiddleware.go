package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements token bucket algorithm for rate limiting
type RateLimiter struct {
	// mu protects the buckets map
	mu sync.RWMutex
	// buckets stores the token bucket for each IP
	buckets map[string]*tokenBucket
	// requestsPerSecond is the rate limit (tokens per second)
	requestsPerSecond float64
	// burstSize is the maximum tokens that can be accumulated
	burstSize int64
	// cleanupInterval is how often to clean old buckets
	cleanupInterval time.Duration
}

// tokenBucket tracks available tokens and last refill time
type tokenBucket struct {
	tokens    float64
	lastRefil time.Time
}

// NewRateLimiter creates a new rate limiter
// requestsPerSecond: max requests per second (e.g., 10.0 for 10 req/sec)
// burstSize: maximum requests allowed in a burst (e.g., 100)
func NewRateLimiter(requestsPerSecond float64, burstSize int64) *RateLimiter {
	rl := &RateLimiter{
		buckets:           make(map[string]*tokenBucket),
		requestsPerSecond: requestsPerSecond,
		burstSize:         burstSize,
		cleanupInterval:   5 * time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanupRoutine()

	return rl
}

// cleanupRoutine periodically removes old buckets to free memory
func (rl *RateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, bucket := range rl.buckets {
			// Remove bucket if it hasn't been used in the last hour
			if now.Sub(bucket.lastRefil) > time.Hour {
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	bucket, exists := rl.buckets[ip]
	if !exists {
		// New bucket for this IP
		bucket = &tokenBucket{
			tokens:    float64(rl.burstSize),
			lastRefil: now,
		}
		rl.buckets[ip] = bucket
	}

	// Calculate tokens to add based on time elapsed
	timePassed := now.Sub(bucket.lastRefil).Seconds()
	tokensToAdd := timePassed * rl.requestsPerSecond

	// Add tokens but don't exceed burst size
	bucket.tokens = min(float64(rl.burstSize), bucket.tokens+tokensToAdd)
	bucket.lastRefil = now

	// Check if we have tokens available
	if bucket.tokens >= 1.0 {
		bucket.tokens--
		return true
	}

	return false
}

// Helper function to return minimum of two values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// RateLimitMiddleware returns middleware that rate limits requests
func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP address
			ip := getClientIP(r)

			// Check if request is allowed
			if !limiter.Allow(ip) {
				w.Header().Set("Retry-After", "60")
				w.Header().Set("X-RateLimit-Limit", "10")
				w.Header().Set("X-RateLimit-Remaining", "0")
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Add rate limit headers to response
			w.Header().Set("X-RateLimit-Limit", "10")
			w.Header().Set("X-RateLimit-Remaining", "9")

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies like Vercel, nginx)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		return xff
	}

	// Check X-Real-IP header (alternative proxy header)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	return r.RemoteAddr
}

// RateLimitHandlerFunc is a convenience function for wrapping http.HandlerFunc
// Usage: mux.Handle("/auth/login", RateLimitHandlerFunc(limiter, 10, 100)(http.HandlerFunc(handler)))
func RateLimitHandlerFunc(requestsPerSecond float64, burstSize int64) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(requestsPerSecond, burstSize)
	return RateLimitMiddleware(limiter)
}
