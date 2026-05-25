// Package linenum tracks absolute and relative line numbers during log
// processing, enabling callers to correlate output lines back to their
// original positions in the source file.
package linenum

import "sync/atomic"

// Tracker counts lines seen and provides helpers for relative offsets.
type Tracker struct {
	total   atomic.Int64
	matched atomic.Int64
}

// New returns an initialised Tracker with all counters at zero.
func New() *Tracker {
	return &Tracker{}
}

// IncTotal records that one more source line has been consumed.
// It returns the new absolute line number (1-based).
func (t *Tracker) IncTotal() int64 {
	return t.total.Add(1)
}

// IncMatched records that one more line passed all filters.
// It returns the new matched-line count (1-based).
func (t *Tracker) IncMatched() int64 {
	return t.matched.Add(1)
}

// Total returns the number of source lines consumed so far.
func (t *Tracker) Total() int64 {
	return t.total.Load()
}

// Matched returns the number of lines that passed all filters.
func (t *Tracker) Matched() int64 {
	return t.matched.Load()
}

// Reset zeroes both counters.
func (t *Tracker) Reset() {
	t.total.Store(0)
	t.matched.Store(0)
}

// Annotation holds per-line numbering metadata.
type Annotation struct {
	Absolute int64 // position in source file (1-based)
	Relative int64 // position among matched lines (1-based)
}

// Annotate calls IncTotal, conditionally calls IncMatched, and returns
// an Annotation. Pass matched=true when the line survived all filters.
func (t *Tracker) Annotate(matched bool) Annotation {
	abs := t.IncTotal()
	var rel int64
	if matched {
		rel = t.IncMatched()
	}
	return Annotation{Absolute: abs, Relative: rel}
}
