package ratelimiter

import (
	"sync"
	"sync/atomic"
	"time"
)

type FixedWindowRateLimiter struct {
	mu        sync.RWMutex
	count     int32
	limit     int32
	window    time.Duration
	timer     *time.Timer
	startTime time.Time
}

func NewFixedWindowRateLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	rl := &FixedWindowRateLimiter{
		limit:     int32(limit),
		window:    window,
		startTime: time.Now(),
	}
	rl.timer = time.AfterFunc(window, rl.resetCount)
	return rl
}

func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if atomic.LoadInt32(&rl.count) < rl.limit {
		atomic.AddInt32(&rl.count, 1)
		return true, 0
	}
	retryAfter := rl.window - time.Since(rl.startTime)
	return false, retryAfter
}

func (rl *FixedWindowRateLimiter) resetCount() {
	atomic.StoreInt32(&rl.count, 0)
	rl.timer.Reset(rl.window)
}
