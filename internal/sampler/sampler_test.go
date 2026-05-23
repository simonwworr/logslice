package sampler

import (
	"testing"
)

func TestNew_InvalidN(t *testing.T) {
	_, err := New(Interval, 0, 0)
	if err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
	_, err = New(Interval, -3, 0)
	if err == nil {
		t.Fatal("expected error for n=-3, got nil")
	}
}

func TestNew_N1_KeepsAll(t *testing.T) {
	s, err := New(Interval, 1, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 20; i++ {
		if !s.Keep() {
			t.Errorf("expected Keep()=true for every line when n=1, failed at call %d", i+1)
		}
	}
}

func TestInterval_EveryNthLine(t *testing.T) {
	const n = 3
	s, err := New(Interval, n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// With interval strategy counter starts at 0, increments before check.
	// Kept lines: counter % n == 0 → counters 3, 6, 9 ...
	kept := 0
	for i := 0; i < 9; i++ {
		if s.Keep() {
			kept++
		}
	}
	if kept != 3 {
		t.Errorf("expected 3 kept lines out of 9 with n=3, got %d", kept)
	}
}

func TestReset_ResetsCounter(t *testing.T) {
	s, _ := New(Interval, 3, 0)
	// Consume 2 lines (not yet at multiple of 3).
	s.Keep()
	s.Keep()
	s.Reset()
	// After reset counter is 0; next Keep() increments to 1 → not kept.
	if s.Keep() {
		t.Error("expected first line after Reset not to be kept (counter=1, n=3)")
	}
}

func TestRandom_ApproximateRate(t *testing.T) {
	const n = 10
	const trials = 10_000
	s, err := New(Random, n, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	kept := 0
	for i := 0; i < trials; i++ {
		if s.Keep() {
			kept++
		}
	}
	// Expect ~1000 ± 200 (within 2 sigma for binomial).
	if kept < 800 || kept > 1200 {
		t.Errorf("random sampling rate out of expected range: kept=%d / %d", kept, trials)
	}
}
