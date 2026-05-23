package context

import (
	"context"
	"time"
)

// Key is an unexported type for context keys in this package.
type Key int

const (
	// deadlineKey is the context key for a processing deadline.
	deadlineKey Key = iota
	// startTimeKey is the context key for the slice start time.
	startTimeKey
	// endTimeKey is the context key for the slice end time.
	endTimeKey
)

// WithDeadline returns a new context that will be cancelled after d duration.
func WithDeadline(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d)
}

// WithTimeRange stores the log slice time range in the context.
func WithTimeRange(ctx context.Context, start, end time.Time) context.Context {
	ctx = context.WithValue(ctx, startTimeKey, start)
	ctx = context.WithValue(ctx, endTimeKey, end)
	return ctx
}

// TimeRange retrieves the start and end times from the context.
// Returns zero times and false if not set.
func TimeRange(ctx context.Context) (start, end time.Time, ok bool) {
	s, sok := ctx.Value(startTimeKey).(time.Time)
	e, eok := ctx.Value(endTimeKey).(time.Time)
	if !sok || !eok {
		return time.Time{}, time.Time{}, false
	}
	return s, e, true
}

// Done returns a channel that is closed when the context is cancelled or times out.
func Done(ctx context.Context) <-chan struct{} {
	return ctx.Done()
}

// Err returns the context error, if any.
func Err(ctx context.Context) error {
	return ctx.Err()
}
