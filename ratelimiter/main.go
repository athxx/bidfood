package main

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	mu        sync.Mutex
	timestamp int64
	tokens    int
	rate      int
}

func NewRateLimiter(rate int) *RateLimiter {
	return &RateLimiter{
		rate: rate,
	}
}

func (l *RateLimiter) Allow() {
	for {
		l.mu.Lock()
		now := time.Now().Unix()
		if now != l.timestamp {
			l.timestamp = now
			l.tokens = 0
		}
		if l.tokens < l.rate {
			l.tokens++
			l.mu.Unlock()
			return
		}
		l.mu.Unlock()
	}
}

// ------------------------------ normal limiter base on token -------------------------
// Limiter struct to control rates
type Limiter struct {
	tokens   chan struct{}
	interval time.Duration
}

// NewLimiter creates a new RateLimiter allow `rate` operations per sec
func NewLimiter(rate int) *Limiter {
	l := &Limiter{
		tokens:   make(chan struct{}, rate),
		interval: time.Second / time.Duration(rate),
	}
	go func() {
		ticker := time.NewTicker(l.interval)
		defer ticker.Stop()
		for range ticker.C {
			select {
			case l.tokens <- struct{}{}:
			default:
			}
		}
	}()

	return l
}

func (l *Limiter) Allow() {
	<-l.tokens
}

func main() {
	limiter := NewRateLimiter(5)

	for i := 0; i < 30; i++ {
		limiter.Allow()
		fmt.Printf("Operation %d at %v\n", i, time.Now().Format("2006-01-02 15:04:05.000000000000000000"))
	}

	fmt.Println("----------------------------------- Split -----------------------------------")
	l := NewLimiter(5)

	for i := 0; i < 30; i++ {
		l.Allow()
		fmt.Printf("Operation %d at %v\n", i, time.Now().Format("2006-01-02 15:04:05.000000000000000000"))
	}
}
