package linenum_test

import (
	"sync"
	"testing"

	"github.com/yourorg/logslice/internal/linenum"
)

func TestNew_ZeroCounters(t *testing.T) {
	tr := linenum.New()
	if tr.Total() != 0 || tr.Matched() != 0 {
		t.Fatalf("expected zero counters, got total=%d matched=%d", tr.Total(), tr.Matched())
	}
}

func TestIncTotal_Increments(t *testing.T) {
	tr := linenum.New()
	if got := tr.IncTotal(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
	tr.IncTotal()
	if tr.Total() != 2 {
		t.Fatalf("expected Total=2, got %d", tr.Total())
	}
}

func TestIncMatched_Increments(t *testing.T) {
	tr := linenum.New()
	tr.IncMatched()
	tr.IncMatched()
	if tr.Matched() != 2 {
		t.Fatalf("expected 2, got %d", tr.Matched())
	}
}

func TestReset_ZerosBoth(t *testing.T) {
	tr := linenum.New()
	tr.IncTotal()
	tr.IncMatched()
	tr.Reset()
	if tr.Total() != 0 || tr.Matched() != 0 {
		t.Fatal("Reset did not zero counters")
	}
}

func TestAnnotate_MatchedLine(t *testing.T) {
	tr := linenum.New()
	a := tr.Annotate(true)
	if a.Absolute != 1 {
		t.Fatalf("expected Absolute=1, got %d", a.Absolute)
	}
	if a.Relative != 1 {
		t.Fatalf("expected Relative=1, got %d", a.Relative)
	}
}

func TestAnnotate_UnmatchedLine(t *testing.T) {
	tr := linenum.New()
	a := tr.Annotate(false)
	if a.Absolute != 1 {
		t.Fatalf("expected Absolute=1, got %d", a.Absolute)
	}
	if a.Relative != 0 {
		t.Fatalf("expected Relative=0 for unmatched, got %d", a.Relative)
	}
	if tr.Matched() != 0 {
		t.Fatal("Matched counter should remain 0")
	}
}

func TestAnnotate_MixedLines(t *testing.T) {
	tr := linenum.New()
	// line 1 – unmatched
	tr.Annotate(false)
	// line 2 – matched
	a := tr.Annotate(true)
	if a.Absolute != 2 || a.Relative != 1 {
		t.Fatalf("unexpected annotation %+v", a)
	}
}

func TestAnnotate_Concurrent(t *testing.T) {
	tr := linenum.New()
	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			tr.Annotate(true)
		}()
	}
	wg.Wait()
	if tr.Total() != goroutines {
		t.Fatalf("expected Total=%d, got %d", goroutines, tr.Total())
	}
	if tr.Matched() != goroutines {
		t.Fatalf("expected Matched=%d, got %d", goroutines, tr.Matched())
	}
}
