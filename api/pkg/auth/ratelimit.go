package auth

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter manages rate limiting for authentication endpoints
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int           // Max requests per window
	window   time.Duration // Time window for rate limiting
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
	
	// Start cleanup goroutine to prevent memory leaks
	go rl.cleanup()
	
	return rl
}

// Allow checks if the request should be allowed based on rate limits
func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-rl.window)
	
	// Get existing requests for this IP
	requests := rl.requests[clientIP]
	
	// Filter out old requests
	var validRequests []time.Time
	for _, req := range requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}
	
	// Check if we're within the limit
	if len(validRequests) >= rl.limit {
		rl.requests[clientIP] = validRequests
		return false
	}
	
	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[clientIP] = validRequests
	
	return true
}

// cleanup removes old entries to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window * 2) // Keep extra buffer
		
		for ip, requests := range rl.requests {
			var validRequests []time.Time
			for _, req := range requests {
				if req.After(cutoff) {
					validRequests = append(validRequests, req)
				}
			}
			
			if len(validRequests) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = validRequests
			}
		}
		rl.mutex.Unlock()
	}
}

// RateLimitMiddleware returns middleware for rate limiting authentication endpoints
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP address
		clientIP := getClientIP(r)
		
		// Check rate limit
		if !rl.Allow(clientIP) {
			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", "5")
			w.Header().Set("X-RateLimit-Window", "300") // 5 minutes in seconds
			w.Header().Set("Retry-After", "300")
			
			http.Error(w, "Too many authentication attempts. Please try again later.", http.StatusTooManyRequests)
			return
		}
		
		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the real client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for load balancers/proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can be a comma-separated list, take the first one
		if len(xff) > 0 {
			return xff
		}
	}
	
	// Check X-Real-IP header (for reverse proxies)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	
	// Fall back to RemoteAddr
	return r.RemoteAddr
}