package linecount

import (
	"sync"
	"testing"
)

func TestNew_ZeroCounters(t *testing.T) {
	c := New()
	if c.Total() != 0 || c.Matched() != 0 || c.Skipped() != 0 {
		t.Fatalf("expected all zero counters, got total=%d matched=%d skipped=%d",
			c.Total(), c.Matched(), c.Skipped())
	}
}

func TestIncTotal(t *testing.T) {
	c := New()
	c.IncTotal()
	c.IncTotal()
	if c.Total() != 2 {
		t.Fatalf("expected 2, got %d", c.Total())
	}
}

func TestIncMatched(t *testing.T) {
	c := New()
	for i := 0; i < 5; i++ {
		c.IncMatched()
	}
	if c.Matched() != 5 {
		t.Fatalf("expected 5, got %d", c.Matched())
	}
}

func TestIncSkipped(t *testing.T) {
	c := New()
	c.IncSkipped()
	if c.Skipped() != 1 {
		t.Fatalf("expected 1, got %d", c.Skipped())
	}
}

func TestReset_ZerosAll(t *testing.T) {
	c := New()
	c.IncTotal()
	c.IncMatched()
	c.IncSkipped()
	c.Reset()
	snap := c.Snap()
	if snap.Total != 0 || snap.Matched != 0 || snap.Skipped != 0 {
		t.Fatalf("expected zeros after reset, got %+v", snap)
	}
}

func TestSnap_CapturesValues(t *testing.T) {
	c := New()
	c.IncTotal()
	c.IncTotal()
	c.IncMatched()
	snap := c.Snap()
	if snap.Total != 2 || snap.Matched != 1 || snap.Skipped != 0 {
		t.Fatalf("unexpected snapshot %+v", snap)
	}
}

func TestConcurrent_SafeIncrements(t *testing.T) {
	c := New()
	const goroutines = 50
	const perGoroutine = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				c.IncTotal()
				c.IncMatched()
			}
		}()
	}
	wg.Wait()
	expected := int64(goroutines * perGoroutine)
	if c.Total() != expected {
		t.Fatalf("total: expected %d, got %d", expected, c.Total())
	}
	if c.Matched() != expected {
		t.Fatalf("matched: expected %d, got %d", expected, c.Matched())
	}
}
