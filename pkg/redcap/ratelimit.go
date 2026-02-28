package redcap

import (
	"context"
	"golang.org/x/time/rate"
)

const (
	DefaultRPS   = 10
	DefaultBurst = 10
)

// RateLimiter defines the contract. 
// Note: Wait must return error to handle context cancellation.
type RateLimiter interface {
	Wait(ctx context.Context) error
	SetRate(rps float64)
}

type limiter struct {
	r *rate.Limiter
}

// NewRateLimiterWithDefaultOpts returns the interface, not a pointer to it.
func NewRateLimiterWithDefaultOpts() RateLimiter {
	return &limiter{
		r: rate.NewLimiter(rate.Limit(DefaultRPS), DefaultBurst),
	}
}

// Wait matches the interface signature exactly.
func (l *limiter) Wait(ctx context.Context) error {
	return l.r.Wait(ctx)
}

// SetRate converts the float64 to rate.Limit internally.
func (l *limiter) SetRate(rps float64) {
	l.r.SetLimit(rate.Limit(rps))
}
