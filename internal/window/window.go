// Package window provides sliding and tumbling window aggregation
// over timestamped log lines, grouping them into fixed-duration buckets.
package window

import (
	"fmt"
	"time"
)

// Bucket holds log lines that fall within a specific time window.
type Bucket struct {
	Start time.Time
	End   time.Time
	Lines []string
	Count int
}

// String returns a human-readable summary of the bucket.
func (b *Bucket) String() string {
	return fmt.Sprintf("[%s – %s] %d lines",
		b.Start.Format(time.RFC3339),
		b.End.Format(time.RFC3339),
		b.Count,
	)
}

// Windower groups timestamped lines into fixed-duration buckets.
type Windower struct {
	size    time.Duration
	buckets []*Bucket
	current *Bucket
}

// New creates a Windower that groups lines into buckets of the given duration.
// Returns an error if size is non-positive.
func New(size time.Duration) (*Windower, error) {
	if size <= 0 {
		return nil, fmt.Errorf("window: size must be positive, got %s", size)
	}
	return &Windower{size: size}, nil
}

// Add places line into the appropriate time bucket based on ts.
func (w *Windower) Add(ts time.Time, line string) {
	if w.current == nil || !ts.Before(w.current.End) {
		w.flush(ts)
	}
	w.current.Lines = append(w.current.Lines, line)
	w.current.Count++
}

// flush closes the current bucket and opens a new one aligned to ts.
func (w *Windower) flush(ts time.Time) {
	start := ts.Truncate(w.size)
	b := &Bucket{
		Start: start,
		End:   start.Add(w.size),
	}
	w.buckets = append(w.buckets, b)
	w.current = b
}

// Buckets returns all buckets collected so far, including the current one.
func (w *Windower) Buckets() []*Bucket {
	return w.buckets
}

// Reset clears all buckets and resets the windower to its initial state.
func (w *Windower) Reset() {
	w.buckets = nil
	w.current = nil
}
