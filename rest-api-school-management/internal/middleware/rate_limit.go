package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type rateLimit struct {
	mu        sync.Mutex
	visitors  map[string]int
	limit     int
	resetTime time.Duration
}

func NewRateLimit(limit int, resetTime time.Duration) *rateLimit {
	rl := &rateLimit{
		visitors:  make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}
	go rl.resetVisitorCount()
	return rl
}

func (rl *rateLimit) resetVisitorCount() {
	for {
		time.Sleep(rl.resetTime)
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}

func (rl *rateLimit) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		ip := r.RemoteAddr // for production IP with some better mechanism here
		rl.visitors[ip]++
		fmt.Printf("Visitors count from %v is %v\n", ip, rl.visitors[ip])
		if rl.visitors[ip] > rl.limit {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
