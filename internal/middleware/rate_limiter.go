package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"external-backend-go/internal/logger"
)

type RateLimiter struct {
	enabled         bool
	rps             float64
	burst           int
	clients         map[string]*rate.Limiter
	mu              sync.Mutex
	cleanupInterval time.Duration
}

func NewRateLimiter(enabled bool, rps float64, burst int, cleanupInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		enabled:         enabled,
		rps:             rps,
		burst:           burst,
		clients:         make(map[string]*rate.Limiter),
		cleanupInterval: cleanupInterval,
	}
	if enabled {
		go rl.cleanupClients()
	}
	return rl
}

func (rl *RateLimiter) Allow(ip string) (bool, time.Duration) {
	if !rl.enabled {
		return true, 0
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.clients[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(rl.rps), rl.burst)
		rl.clients[ip] = limiter
	}

	if !limiter.Allow() {
		return false, time.Duration(float64(time.Second) / rl.rps)
	}
	return true, 0
}

func (rl *RateLimiter) cleanupClients() {
	for range time.Tick(rl.cleanupInterval) {
		rl.mu.Lock()
		for ip, limiter := range rl.clients {
			if limiter.Tokens() >= float64(rl.burst) {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
		logger.Info("Rate limiter client map cleaned up. Current clients: %d", len(rl.clients))
	}
}

func RateLimiterMiddleware(limiter *RateLimiter, rateLimitExceededErr func(w http.ResponseWriter, r *http.Request, retryAfter string)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.enabled {
				next.ServeHTTP(w, r)
				return
			}

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			allow, retryAfter := limiter.Allow(ip)
			if !allow {
				rateLimitExceededErr(w, r, retryAfter.String())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
