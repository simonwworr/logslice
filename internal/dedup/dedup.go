// Package dedup provides line deduplication for log output.
// It tracks seen lines using a configurable hash-based window
// and filters out consecutive or global duplicates.
package dedup

import (
	"crypto/sha256"
	"encoding/hex"
)

// Mode controls how deduplication is applied.
type Mode int

const (
	// Consecutive removes only back-to-back duplicate lines.
	Consecutive Mode = iota
	// Global removes any line seen previously in the stream.
	Global
)

// Deduplicator filters duplicate log lines.
type Deduplicator struct {
	mode    Mode
	seen    map[string]struct{}
	lastKey string
}

// New creates a Deduplicator with the given mode.
func New(mode Mode) *Deduplicator {
	return &Deduplicator{
		mode: mode,
		seen: make(map[string]struct{}),
	}
}

// IsDuplicate reports whether line has been seen according to the mode.
// It returns true if the line should be suppressed.
func (d *Deduplicator) IsDuplicate(line string) bool {
	key := hash(line)
	switch d.mode {
	case Consecutive:
		if key == d.lastKey {
			return true
		}
		d.lastKey = key
		return false
	case Global:
		if _, ok := d.seen[key]; ok {
			return true
		}
		d.seen[key] = struct{}{}
		d.lastKey = key
		return false
	}
	return false
}

// Reset clears all deduplication state.
func (d *Deduplicator) Reset() {
	d.seen = make(map[string]struct{})
	d.lastKey = ""
}

// UniqueCount returns the number of unique lines seen (Global mode only).
func (d *Deduplicator) UniqueCount() int {
	return len(d.seen)
}

func hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
