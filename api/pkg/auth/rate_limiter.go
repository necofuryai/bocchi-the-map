package auth

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/necofuryai/bocchi-the-map/api/pkg/errors"
	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
	"github.com/necofuryai/bocchi-the-map/api/pkg/monitoring"
)

// RateLimiter provides rate limiting functionality for authentication endpoints
type RateLimiter struct {
	requests     map[string]*requestCounter
	mutex        sync.RWMutex
	limit        int
	window       time.Duration
	cleanupTimer *time.Ticker
}

// requestCounter tracks requests for a specific IP address
type requestCounter struct {
	count     int
	firstSeen time.Time
	lastSeen  time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests:     make(map[string]*requestCounter),
		limit:        limit,
		window:       window,
		cleanupTimer: time.NewTicker(window / 2), // Cleanup every half window
	}

	// Start cleanup goroutine
	go rl.cleanupRoutine()

	logger.InfoWithFields("Rate limiter initialized", map[string]interface{}{
		"limit":  limit,
		"window": window.String(),
	})

	return rl
}

// IsAllowed checks if a request from the given IP is allowed
func (rl *RateLimiter) IsAllowed(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	
	// Get or create counter for this IP
	counter, exists := rl.requests[ip]
	if !exists {
		rl.requests[ip] = &requestCounter{
			count:     1,
			firstSeen: now,
			lastSeen:  now,
		}
		return true
	}

	// Check if the window has expired
	if now.Sub(counter.firstSeen) > rl.window {
		// Reset the counter for a new window
		counter.count = 1
		counter.firstSeen = now
		counter.lastSeen = now
		return true
	}

	// Update last seen time
	counter.lastSeen = now

	// Check if limit is exceeded
	if counter.count >= rl.limit {
		logger.InfoWithFields("Rate limit exceeded", map[string]interface{}{
			"ip":    ip,
			"count": counter.count,
			"limit": rl.limit,
		})
		return false
	}

	// Increment counter
	counter.count++
	return true
}

// GetRemainingRequests returns the number of remaining requests for an IP
func (rl *RateLimiter) GetRemainingRequests(ip string) int {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	counter, exists := rl.requests[ip]
	if !exists {
		return rl.limit
	}

	// Check if window has expired
	if time.Since(counter.firstSeen) > rl.window {
		return rl.limit
	}

	remaining := rl.limit - counter.count
	if remaining < 0 {
		return 0
	}

	return remaining
}

// GetResetTime returns when the rate limit will reset for an IP
func (rl *RateLimiter) GetResetTime(ip string) time.Time {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	counter, exists := rl.requests[ip]
	if !exists {
		return time.Now()
	}

	return counter.firstSeen.Add(rl.window)
}

// Middleware returns a Chi middleware function for rate limiting
func (rl *RateLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP
			ip := rl.getClientIP(r)
			
			// Add monitoring context
			ctx := monitoring.StartTrace(r.Context(), "auth.rate_limit_check")
			defer monitoring.EndTrace(ctx)

			// Check rate limit
			if !rl.IsAllowed(ip) {
				rl.handleRateLimitExceeded(w, r, ip)
				return
			}

			// Add rate limit headers
			remaining := rl.GetRemainingRequests(ip)
			resetTime := rl.GetResetTime(ip)
			
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the real client IP from the request
func (rl *RateLimiter) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (most common for proxies)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Check X-Real-IP header (used by nginx)
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		if net.ParseIP(xri) != nil {
			return xri
		}
	}

	// Check CF-Connecting-IP header (used by Cloudflare)
	cfip := r.Header.Get("CF-Connecting-IP")
	if cfip != "" {
		if net.ParseIP(cfip) != nil {
			return cfip
		}
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If there's no port, the entire string might be an IP
		if net.ParseIP(r.RemoteAddr) != nil {
			return r.RemoteAddr
		}
		// Return a default IP if parsing fails
		return "unknown"
	}

	return ip
}

// handleRateLimitExceeded handles rate limit exceeded responses
func (rl *RateLimiter) handleRateLimitExceeded(w http.ResponseWriter, r *http.Request, ip string) {
	// Log rate limit exceeded
	logger.InfoWithFields("Rate limit exceeded", map[string]interface{}{
		"ip":     ip,
		"path":   r.URL.Path,
		"method": r.Method,
		"limit":  rl.limit,
		"window": rl.window.String(),
	})

	// Add monitoring metrics
	monitoring.RecordRateLimitExceeded(r.Context(), ip)

	// Create domain error
	domainErr := errors.New(errors.ErrTypeRateLimit, "Rate limit exceeded. Please try again later.")

	// Set response headers
	remaining := rl.GetRemainingRequests(ip)
	resetTime := rl.GetResetTime(ip)
	retryAfter := int(time.Until(resetTime).Seconds())
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
	w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
	w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
	w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))

	// Set status code
	w.WriteHeader(domainErr.ToHTTPStatus())

	// Write error response
	errorResponse := fmt.Sprintf(`{
		"error": {
			"type": "%s",
			"message": "%s",
			"retry_after": %d
		},
		"request_id": "%s"
	}`, domainErr.Type, domainErr.Message, retryAfter, monitoring.GetRequestID(r.Context()))

	w.Write([]byte(errorResponse))
}

// cleanupRoutine periodically cleans up expired entries
func (rl *RateLimiter) cleanupRoutine() {
	for range rl.cleanupTimer.C {
		rl.cleanupExpiredEntries()
	}
}

// cleanupExpiredEntries removes expired entries from the rate limiter
func (rl *RateLimiter) cleanupExpiredEntries() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	var toDelete []string

	for ip, counter := range rl.requests {
		// Remove entries that haven't been seen for longer than the window
		if now.Sub(counter.lastSeen) > rl.window {
			toDelete = append(toDelete, ip)
		}
	}

	// Delete expired entries
	for _, ip := range toDelete {
		delete(rl.requests, ip)
	}

	if len(toDelete) > 0 {
		logger.InfoWithFields("Rate limiter cleanup completed", map[string]interface{}{
			"cleaned_entries": len(toDelete),
			"remaining":       len(rl.requests),
		})
	}
}

// Stop stops the rate limiter and cleanup routines
func (rl *RateLimiter) Stop() {
	if rl.cleanupTimer != nil {
		rl.cleanupTimer.Stop()
	}
	
	logger.Info("Rate limiter stopped")
}

// GetStats returns statistics about the rate limiter
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	return map[string]interface{}{
		"active_ips": len(rl.requests),
		"limit":      rl.limit,
		"window":     rl.window.String(),
	}
}