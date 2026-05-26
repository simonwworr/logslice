// Package linecount provides a streaming line counter that tracks total,
// matched, and skipped lines with optional periodic progress callbacks.
package linecount

import "sync/atomic"

// Counter tracks line processing statistics using atomic operations
// so it is safe for concurrent use.
type Counter struct {
	total   atomic.Int64
	matched atomic.Int64
	skipped atomic.Int64
}

// New returns an initialised Counter ready for use.
func New() *Counter {
	return &Counter{}
}

// IncTotal records one additional line seen by the pipeline.
func (c *Counter) IncTotal() {
	c.total.Add(1)
}

// IncMatched records one additional line that passed all filters.
func (c *Counter) IncMatched() {
	c.matched.Add(1)
}

// IncSkipped records one additional line that was dropped.
func (c *Counter) IncSkipped() {
	c.skipped.Add(1)
}

// Total returns the number of lines seen.
func (c *Counter) Total() int64 {
	return c.total.Load()
}

// Matched returns the number of lines that passed all filters.
func (c *Counter) Matched() int64 {
	return c.matched.Load()
}

// Skipped returns the number of lines that were dropped.
func (c *Counter) Skipped() int64 {
	return c.skipped.Load()
}

// Reset zeroes all counters.
func (c *Counter) Reset() {
	c.total.Store(0)
	c.matched.Store(0)
	c.skipped.Store(0)
}

// Snapshot returns a point-in-time copy of the counters as plain integers.
type Snapshot struct {
	Total   int64
	Matched int64
	Skipped int64
}

// Snap captures the current counter values atomically.
func (c *Counter) Snap() Snapshot {
	return Snapshot{
		Total:   c.total.Load(),
		Matched: c.matched.Load(),
		Skipped: c.skipped.Load(),
	}
}
