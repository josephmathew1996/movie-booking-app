package ratelimiter

import (
	"sync"
	"sync/atomic"
	"time"
)

type FixedWindowClientsRateLimiter struct {
	mu      sync.RWMutex
	clients map[string]*clientInfo
	limit   int
	window  time.Duration
}

type clientInfo struct {
	count int32
	start time.Time
	timer *time.Timer
}

func NewFixedWindowClientsRateLimiter(limit int, window time.Duration) *FixedWindowClientsRateLimiter {
	return &FixedWindowClientsRateLimiter{
		clients: make(map[string]*clientInfo),
		limit:   limit,
		window:  window,
	}
}

func (rl *FixedWindowClientsRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	client, exists := rl.clients[ip]
	//if a client does not exist, create a new client info
	if !exists {
		client = &clientInfo{start: time.Now()}
		// set the timer to reset the count for each client after a `window` duration
		client.timer = time.AfterFunc(rl.window, func() {
			rl.resetCount(ip)
		})
		rl.clients[ip] = client
	}
	if atomic.LoadInt32(&client.count) < int32(rl.limit) {
		atomic.AddInt32(&client.count, 1)
		return true, 0
	}
	retryAfter := rl.window - time.Since(client.start)
	return false, retryAfter
}

func (rl *FixedWindowClientsRateLimiter) resetCount(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.clients, ip)
}
