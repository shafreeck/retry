package retry

import (
	"context"
	"time"
)

// Retriable is an error type which can be retried
type Retriable error

// Retry ensures that the do function will be executed until some condition being satisfied
type Retry interface {
	Ensure(ctx context.Context, do func() error) error
}

type retry struct {
	backoff BackoffStrategy
	base    time.Duration
}

var r = New()

// Ensure keeps retring until ctx is done
func (r *retry) Ensure(ctx context.Context, do func() error) error {
	duration := r.base
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := do(); err != nil {
			if _, ok := err.(Retriable); ok {
				if r.backoff != nil {
					duration = r.backoff(duration)

					time.Sleep(duration)
				}
				continue
			}
			return err
		}

		return nil
	}
}

// Option is an option to new a Retry object
type Option func(r *retry)

// BackoffStrategy defines the backoff strategy of retry
type BackoffStrategy func(last time.Duration) time.Duration

// WithBackoff replace the default backoff function
func WithBackoff(backoff BackoffStrategy) Option {
	return func(r *retry) {
		r.backoff = backoff
	}
}

// WithBase set the first delay duration, default 10ms
func WithBaseDelay(base time.Duration) Option {
	return func(r *retry) {
		r.base = base
	}
}

// New a retry object
func New(opts ...Option) Retry {
	r := &retry{base: 10 * time.Millisecond, backoff: Exponential(2)}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Exponential generates backoff duration by expoential
func Exponential(factor float64) BackoffStrategy {
	return func(last time.Duration) time.Duration {
		return time.Duration(float64(last) * factor)
	}
}

// Ensure keeps retring until ctx is done, it use a default retry object
func Ensure(ctx context.Context, do func() error) error {
	return r.Ensure(ctx, do)
}
