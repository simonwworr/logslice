// Package limiter provides line-count and byte-count limiting for log output.
package limiter

import "errors"

// ErrLimitReached is returned when the configured limit has been reached.
var ErrLimitReached = errors.New("limit reached")

// Limiter tracks how many lines and bytes have been processed and signals
// when a configured threshold has been exceeded.
type Limiter struct {
	maxLines int64
	maxBytes int64
	lines    int64
	bytes    int64
}

// New creates a Limiter. Use 0 to disable a limit.
func New(maxLines, maxBytes int64) (*Limiter, error) {
	if maxLines < 0 {
		return nil, errors.New("maxLines must be >= 0")
	}
	if maxBytes < 0 {
		return nil, errors.New("maxBytes must be >= 0")
	}
	return &Limiter{maxLines: maxLines, maxBytes: maxBytes}, nil
}

// Add records a single line of the given byte length.
// It returns ErrLimitReached if either limit is now exceeded.
func (l *Limiter) Add(lineLen int64) error {
	l.lines++
	l.bytes += lineLen
	if l.maxLines > 0 && l.lines > l.maxLines {
		return ErrLimitReached
	}
	if l.maxBytes > 0 && l.bytes > l.maxBytes {
		return ErrLimitReached
	}
	return nil
}

// Lines returns the number of lines recorded so far.
func (l *Limiter) Lines() int64 { return l.lines }

// Bytes returns the total bytes recorded so far.
func (l *Limiter) Bytes() int64 { return l.bytes }

// Reset clears all counters without changing the configured limits.
func (l *Limiter) Reset() {
	l.lines = 0
	l.bytes = 0
}
