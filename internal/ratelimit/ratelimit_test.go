package ratelimit

import (
	"testing"
	"time"
)

func TestNew_InvalidRate(t *testing.T) {
	_, err := New(0, 10)
	if err != ErrInvalidRate {
		t.Fatalf("expected ErrInvalidRate, got %v", err)
	}
}

func TestNew_InvalidBurst(t *testing.T) {
	_, err := New(10, 0)
	if err != ErrInvalidBurst {
		t.Fatalf("expected ErrInvalidBurst, got %v", err)
	}
}

func TestAllow_ConsumesTokens(t *testing.T) {
	l, err := New(10, 3)
	if err != nil {
		t.Fatal(err)
	}
	// Burst is 3; first three calls should succeed.
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("call %d: expected Allow()=true", i)
		}
	}
	// Fourth call should be denied (no tokens left).
	if l.Allow() {
		t.Fatal("expected Allow()=false after burst exhausted")
	}
}

func TestAllow_TokensRefillOverTime(t *testing.T) {
	l, err := New(100, 1) // 100 tokens/s, burst 1
	if err != nil {
		t.Fatal(err)
	}
	// Drain the single token.
	l.Allow()

	// Advance the internal clock by 50 ms → +5 tokens, capped at 1.
	base := time.Now()
	l.clock = func() time.Time { return base.Add(50 * time.Millisecond) }

	if !l.Allow() {
		t.Fatal("expected Allow()=true after time advance")
	}
}

func TestReset_RestoresTokens(t *testing.T) {
	l, err := New(1, 5)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		l.Allow()
	}
	if l.Allow() {
		t.Fatal("expected Allow()=false after burst exhausted")
	}
	l.Reset()
	if !l.Allow() {
		t.Fatal("expected Allow()=true after Reset")
	}
}

func TestRate_And_Burst_Accessors(t *testing.T) {
	l, _ := New(42.5, 7)
	if l.Rate() != 42.5 {
		t.Fatalf("Rate: want 42.5, got %v", l.Rate())
	}
	if l.Burst() != 7 {
		t.Fatalf("Burst: want 7, got %v", l.Burst())
	}
}
