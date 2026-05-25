// Package aggregate provides time-bucketed line counting and byte aggregation
// over a sliding or fixed window of log entries.
package aggregate

import (
	"sync"
	"time"
)

// Bucket holds aggregated metrics for a single time window.
type Bucket struct {
	Start time.Time
	End   time.Time
	Lines int64
	Bytes int64
}

// Aggregator accumulates log lines into fixed-duration time buckets.
type Aggregator struct {
	mu       sync.Mutex
	buckets  []Bucket
	duration time.Duration
	current  *Bucket
}

// New creates an Aggregator that groups entries into buckets of the given
// duration. duration must be positive.
func New(duration time.Duration) (*Aggregator, error) {
	if duration <= 0 {
		return nil, errInvalidDuration
	}
	return &Aggregator{duration: duration}, nil
}

// Add records a single log line with its timestamp and byte length.
// Lines with a zero timestamp are ignored.
func (a *Aggregator) Add(ts time.Time, lineBytes int) {
	if ts.IsZero() {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.current == nil || !ts.Before(a.current.End) {
		a.flush()
		start := ts.Truncate(a.duration)
		a.current = &Bucket{
			Start: start,
			End:   start.Add(a.duration),
		}
	}
	a.current.Lines++
	a.current.Bytes += int64(lineBytes)
}

// Buckets returns a snapshot of all completed buckets plus the current
// in-progress bucket (if any). The returned slice is a copy.
func (a *Aggregator) Buckets() []Bucket {
	a.mu.Lock()
	defer a.mu.Unlock()

	out := make([]Bucket, len(a.buckets))
	copy(out, a.buckets)
	if a.current != nil {
		out = append(out, *a.current)
	}
	return out
}

// Reset clears all accumulated state.
func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.buckets = a.buckets[:0]
	a.current = nil
}

// flush moves the current bucket into the completed list.
func (a *Aggregator) flush() {
	if a.current != nil {
		a.buckets = append(a.buckets, *a.current)
		a.current = nil
	}
}
