// Package ratelimit provides a token-bucket rate limiter for controlling
// the throughput of log lines emitted by the pipeline.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter is a token-bucket rate limiter.
type Limiter struct {
	mu       sync.Mutex
	rate     float64   // tokens per second
	burst    float64   // maximum token capacity
	tokens   float64   // current token count
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows up to rate events per second with a
// burst capacity of burst. rate must be positive; burst must be >= 1.
func New(rate float64, burst int) (*Limiter, error) {
	if rate <= 0 {
		return nil, ErrInvalidRate
	}
	if burst < 1 {
		return nil, ErrInvalidBurst
	}
	return &Limiter{
		rate:     rate,
		burst:    float64(burst),
		tokens:   float64(burst),
		lastTick: time.Now(),
		clock:    time.Now,
	}, nil
}

// Allow reports whether one event may proceed at the current time.
// It consumes one token if available.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.burst {
		l.tokens = l.burst
	}

	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// Reset restores the limiter to a full burst capacity as of now.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tokens = l.burst
	l.lastTick = l.clock()
}

// Rate returns the configured rate in events per second.
func (l *Limiter) Rate() float64 { return l.rate }

// Burst returns the configured burst capacity.
func (l *Limiter) Burst() int { return int(l.burst) }
