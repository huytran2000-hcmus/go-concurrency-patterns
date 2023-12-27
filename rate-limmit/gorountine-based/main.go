package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	done := make(chan struct{})
	limiter := NewRateLimiter(2, 2, done)
	var wg sync.WaitGroup

	limiter.Run()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			limiter.Wait()
			fmt.Printf("request no.%d is proccessing\n", idx)
		}(i)
	}

	wg.Wait()
	limiter.Close()
}

type RateLimiter struct {
	buckets chan struct{}
	done    chan struct{}
	ticker  *time.Ticker
}

func NewRateLimiter(limit, rate int, done chan struct{}) *RateLimiter {
	var limiter RateLimiter
	limiter.buckets = make(chan struct{}, limit)
	for i := 0; i < limit; i++ {
		limiter.buckets <- struct{}{}
	}

	limiter.done = done
	limiter.ticker = time.NewTicker(time.Duration(1 / float64(rate) * float64(time.Second)))

	return &limiter
}

func (r *RateLimiter) Run() {
	go func() {
		for {
			select {
			case <-r.ticker.C:
				r.buckets <- struct{}{}
			case <-r.done:
				select {
				case r.buckets <- struct{}{}:
				default:
				}
			}
		}
	}()
}

func (r *RateLimiter) Wait() {
	<-r.buckets
}

func (r *RateLimiter) Close() {
	close(r.done)
	r.ticker.Stop()
}
