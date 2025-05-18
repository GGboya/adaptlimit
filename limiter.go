package adaptlimit

import (
	"context"
	"sync/atomic"
	"time"
)

// Limiter represents a concurrency limiter that uses AIMD algorithm
type Limiter struct {
	limit          atomic.Int64 // current limit
	inflight       atomic.Int64 // current number of inflight requests
	minLimit       int64        // minimum limit
	maxLimit       int64        // maximum limit
	increaseAmount float64      // amount to increase limit by
	decreaseRatio  float64      // ratio to decrease limit by
}

// Option is a function type to modify Limiter
type Option func(*Limiter)

// NewLimiter creates a new Limiter with the given initial limit
func NewLimiter(initialLimit int64, opts ...Option) *Limiter {
	l := &Limiter{
		minLimit:       1,
		maxLimit:       1000,
		increaseAmount: 1,
		decreaseRatio:  0.9,
	}

	// Apply options
	for _, opt := range opts {
		opt(l)
	}

	l.limit.Store(initialLimit)
	return l
}

// WithMinLimit sets the minimum limit
func WithMinLimit(min int64) Option {
	return func(l *Limiter) {
		l.minLimit = min
	}
}

// WithMaxLimit sets the maximum limit
func WithMaxLimit(max int64) Option {
	return func(l *Limiter) {
		l.maxLimit = max
	}
}

// WithIncreaseAmount sets the increase amount
func WithIncreaseAmount(amount float64) Option {
	return func(l *Limiter) {
		l.increaseAmount = amount
	}
}

// WithDecreaseRatio sets the decrease ratio
func WithDecreaseRatio(ratio float64) Option {
	return func(l *Limiter) {
		l.decreaseRatio = ratio
	}
}

// Acquire attempts to acquire a permit from the limiter
func (l *Limiter) Acquire(ctx context.Context) (func(success bool), error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			currentLimit := l.limit.Load()
			currentInflight := l.inflight.Load()

			if currentInflight >= currentLimit {
				time.Sleep(time.Millisecond * 5) // backoff
				continue
			}

			if l.inflight.CompareAndSwap(currentInflight, currentInflight+1) {
				return func(success bool) {
					l.release(success)
				}, nil
			}
		}
	}
}

// release releases a permit back to the limiter
func (l *Limiter) release(success bool) {
	l.inflight.Add(-1)

	if success {
		// Additive increase
		currentLimit := l.limit.Load()
		newLimit := currentLimit + int64(l.increaseAmount)
		if newLimit > l.maxLimit {
			newLimit = l.maxLimit
		}
		l.limit.Store(newLimit)
	} else {
		// Multiplicative decrease
		currentLimit := float64(l.limit.Load())
		newLimit := int64(currentLimit * l.decreaseRatio)
		if newLimit < l.minLimit {
			newLimit = l.minLimit
		}
		l.limit.Store(newLimit)
	}
}

// GetLimit returns the current limit
func (l *Limiter) GetLimit() int64 {
	return l.limit.Load()
}

// GetInflight returns the current number of inflight requests
func (l *Limiter) GetInflight() int64 {
	return l.inflight.Load()
}
