// Package sampler provides line-rate sampling for large log slices.
// It supports both fixed interval (every Nth line) and random reservoir
// sampling strategies, allowing users to reduce output volume while
// preserving a representative view of the log data.
package sampler

import (
	"errors"
	"math/rand"
)

// Strategy defines the sampling approach.
type Strategy int

const (
	// Interval keeps every Nth line (deterministic).
	Interval Strategy = iota
	// Random keeps each line with probability 1/N (non-deterministic).
	Random
)

// Sampler decides whether a given line should be included in output.
type Sampler struct {
	strategy Strategy
	n        int
	counter  int
	rng      *rand.Rand
}

// New creates a Sampler with the given strategy and factor n.
// n must be >= 1; n == 1 means keep every line (no-op sampling).
func New(strategy Strategy, n int, seed int64) (*Sampler, error) {
	if n < 1 {
		return nil, errors.New("sampler: n must be >= 1")
	}
	return &Sampler{
		strategy: strategy,
		n:        n,
		rng:      rand.New(rand.NewSource(seed)), //nolint:gosec
	}, nil
}

// Keep reports whether the current line should be included.
// It must be called exactly once per candidate line.
func (s *Sampler) Keep() bool {
	s.counter++
	switch s.strategy {
	case Interval:
		return s.counter%s.n == 0
	case Random:
		return s.rng.Intn(s.n) == 0
	default:
		return true
	}
}

// Reset resets the internal counter (useful between independent slices).
func (s *Sampler) Reset() {
	s.counter = 0
}
