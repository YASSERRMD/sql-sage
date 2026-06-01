package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/YASSERRMD/sql-sage/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type ipBucket struct {
	count    int
	resetAt  time.Time
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*ipBucket
	limit    int
	window   time.Duration
	idleEvict time.Duration
}

func NewRateLimiter(limitPerMin int) *RateLimiter {
	rl := &RateLimiter{
		buckets:   make(map[string]*ipBucket),
		limit:     limitPerMin,
		window:    time.Minute,
		idleEvict: 5 * time.Minute,
	}
	go rl.gc()
	return rl
}

func (r *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		r.mu.Lock()
		b, ok := r.buckets[ip]
		if !ok || now.After(b.resetAt) {
			b = &ipBucket{count: 0, resetAt: now.Add(r.window)}
			r.buckets[ip] = b
		}
		b.count++
		b.lastSeen = now
		over := b.count > r.limit
		r.mu.Unlock()
		if over {
			c.Header("Retry-After", "60")
			httpx.AbortWithError(c, http.StatusTooManyRequests, "RATE_LIMIT", "too many requests", nil)
			return
		}
		c.Next()
	}
}

func (r *RateLimiter) gc() {
	t := time.NewTicker(time.Minute)
	defer t.Stop()
	for range t.C {
		cut := time.Now().Add(-r.idleEvict)
		r.mu.Lock()
		for k, b := range r.buckets {
			if b.lastSeen.Before(cut) {
				delete(r.buckets, k)
			}
		}
		r.mu.Unlock()
	}
}
