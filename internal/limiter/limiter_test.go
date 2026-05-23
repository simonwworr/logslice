package limiter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/limiter"
)

func TestNew_NegativeLimits(t *testing.T) {
	if _, err := limiter.New(-1, 0); err == nil {
		t.Fatal("expected error for negative maxLines")
	}
	if _, err := limiter.New(0, -1); err == nil {
		t.Fatal("expected error for negative maxBytes")
	}
}

func TestAdd_NoLimit_NeverErrors(t *testing.T) {
	l, _ := limiter.New(0, 0)
	for i := 0; i < 10000; i++ {
		if err := l.Add(128); err != nil {
			t.Fatalf("unexpected error at iteration %d: %v", i, err)
		}
	}
}

func TestAdd_LineLimit(t *testing.T) {
	l, _ := limiter.New(3, 0)
	for i := 0; i < 3; i++ {
		if err := l.Add(10); err != nil {
			t.Fatalf("unexpected error on line %d", i+1)
		}
	}
	if err := l.Add(10); err != limiter.ErrLimitReached {
		t.Fatalf("expected ErrLimitReached, got %v", err)
	}
}

func TestAdd_ByteLimit(t *testing.T) {
	l, _ := limiter.New(0, 50)
	if err := l.Add(30); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := l.Add(25); err != limiter.ErrLimitReached {
		t.Fatalf("expected ErrLimitReached, got %v", err)
	}
}

func TestCounters(t *testing.T) {
	l, _ := limiter.New(0, 0)
	_ = l.Add(10)
	_ = l.Add(20)
	if l.Lines() != 2 {
		t.Errorf("want 2 lines, got %d", l.Lines())
	}
	if l.Bytes() != 30 {
		t.Errorf("want 30 bytes, got %d", l.Bytes())
	}
}

func TestReset(t *testing.T) {
	l, _ := limiter.New(2, 100)
	_ = l.Add(10)
	l.Reset()
	if l.Lines() != 0 || l.Bytes() != 0 {
		t.Error("counters should be zero after Reset")
	}
	// Should not hit limit immediately after reset
	if err := l.Add(10); err != nil {
		t.Fatalf("unexpected error after reset: %v", err)
	}
}
