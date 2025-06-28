package auth

import (
	"fmt"
	"net"
	"net/http"
	"strings"
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

// GetLimit returns the current rate limit value
func (rl *RateLimiter) GetLimit() int {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()
	return rl.limit
}

// GetWindow returns the current rate limit window duration in seconds
func (rl *RateLimiter) GetWindow() int {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()
	return int(rl.window.Seconds())
}

// cleanup removes old entries to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	// Use the rate-limit window (min 1m) as cleanup interval
	cleanupInterval := rl.window
	if cleanupInterval < time.Minute {
		cleanupInterval = time.Minute
	}
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		// Reduce buffer from 2×window to window + 25%
		cutoff := now.Add(-rl.window - (rl.window / 4)) // 25% buffer
		
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
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
			w.Header().Set("X-RateLimit-Window", fmt.Sprintf("%.0f", rl.window.Seconds()))
			w.Header().Set("Retry-After", fmt.Sprintf("%.0f", rl.window.Seconds()))
			
			http.Error(w, "Too many authentication attempts. Please try again later.", http.StatusTooManyRequests)
			return
		}
		
		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

// Trusted proxy IP ranges for validating X-Forwarded-For header
var trustedProxies = []string{
	"130.211.0.0/22",   // Google Cloud Load Balancer
	"35.191.0.0/16",    // Google Cloud Load Balancer
	"34.96.0.0/14",     // Google Cloud Run
	"169.254.169.254/32", // Google Cloud metadata server
	"127.0.0.1/32",     // localhost (for local development)
	"::1/128",          // IPv6 localhost (for local development)
}

// isFromTrustedProxy checks if the request comes from a trusted proxy
func isFromTrustedProxy(remoteAddr string) bool {
	// Extract IP from remoteAddr (format is typically "IP:port")
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// If no port, use the address as is
		host = remoteAddr
	}
	
	clientIP := net.ParseIP(host)
	if clientIP == nil {
		return false
	}
	
	// Check against trusted proxy ranges
	for _, cidr := range trustedProxies {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(clientIP) {
			return true
		}
	}
	
	return false
}

// getClientIP extracts the real client IP from the request with security validation
func getClientIP(r *http.Request) string {
	// Only trust X-Forwarded-For if request comes from a trusted proxy
	if isFromTrustedProxy(r.RemoteAddr) {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			// X-Forwarded-For can be a comma-separated list, take the first one
			ips := strings.Split(xff, ",")
			if len(ips) > 0 {
				// Clean up whitespace and return the first IP
				clientIP := strings.TrimSpace(ips[0])
				if parsedIP := net.ParseIP(clientIP); parsedIP != nil && !parsedIP.IsLoopback() && !parsedIP.IsPrivate() {
					return clientIP
				}
			}
		}
		
		// Check X-Real-IP header as fallback for trusted proxies
		if xri := r.Header.Get("X-Real-IP"); xri != "" {
			if parsedIP := net.ParseIP(xri); parsedIP != nil && !parsedIP.IsLoopback() && !parsedIP.IsPrivate() {
				return xri
			}
		}
	}
	
	// Extract IP from RemoteAddr (format is typically "IP:port")
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If no port, use the address as is
		return r.RemoteAddr
	}
	
	return host
}