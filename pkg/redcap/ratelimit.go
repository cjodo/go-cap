package redcap

import "context"

type RateLimiter interface {
	Wait(ctx context.Context)
	SetRate(rps float64)
}
