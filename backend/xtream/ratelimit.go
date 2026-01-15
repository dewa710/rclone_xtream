package xtream

import (
	"context"
	"time"
)

type RateLimiter struct {
	tokens chan struct{}
}

func NewRateLimiter(rps int) *RateLimiter {
	if rps <= 0 {
		return nil
	}

	rl := &RateLimiter{
		tokens: make(chan struct{}, rps),
	}

	// initial fill
	for i := 0; i < rps; i++ {
		rl.tokens <- struct{}{}
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			for i := len(rl.tokens); i < cap(rl.tokens); i++ {
				rl.tokens <- struct{}{}
			}
		}
	}()

	return rl
}

func (r *RateLimiter) Wait(ctx context.Context) error {
	if r == nil {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.tokens:
		return nil
	}
}
