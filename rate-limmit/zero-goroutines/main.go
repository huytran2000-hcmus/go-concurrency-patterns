package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	limiter := NewEfficientRateLimmiter(2, 2)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			limiter.Wait()
			fmt.Printf("request no.%d is proccessing\n", idx)
		}(i)
	}

	wg.Wait()
}

type EfficientRateLimiter struct {
	nTokens  int
	mu       sync.Mutex
	lastTime time.Time
	limit    int
	period   time.Duration
}

func NewEfficientRateLimmiter(limit, rate int) *EfficientRateLimiter {
	var limiter EfficientRateLimiter
	limiter.nTokens = limit
	limiter.lastTime = time.Now()
	limiter.period = time.Second / time.Duration(rate)
	limiter.limit = limit
	fmt.Printf("period: %v\n", limiter.period)

	return &limiter
}

func (r *EfficientRateLimiter) Wait() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.nTokens > 0 {
		r.nTokens--
		return
	}

	now := time.Now()
	elapsed := now.Sub(r.lastTime)
	newTokens := int(elapsed / r.period)
	if newTokens > 0 {
		r.nTokens += newTokens
		if r.nTokens > r.limit {
			r.nTokens = r.limit
		}
		r.lastTime = now
	}

	if r.nTokens > 0 {
		r.nTokens--
		return
	}

	next := r.lastTime.Add(r.period)
	wait := next.Sub(now)
	<-time.After(wait)
	r.lastTime = next
}
