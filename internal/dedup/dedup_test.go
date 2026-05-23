package dedup

import (
	"testing"
)

func TestConsecutive_NoDuplicates(t *testing.T) {
	d := New(Consecutive)
	lines := []string{"alpha", "beta", "gamma"}
	for _, l := range lines {
		if d.IsDuplicate(l) {
			t.Errorf("expected %q not to be duplicate", l)
		}
	}
}

func TestConsecutive_BackToBack(t *testing.T) {
	d := New(Consecutive)
	if d.IsDuplicate("line") {
		t.Fatal("first occurrence should not be duplicate")
	}
	if !d.IsDuplicate("line") {
		t.Fatal("second consecutive occurrence should be duplicate")
	}
	// Different line resets consecutive tracking
	if d.IsDuplicate("other") {
		t.Fatal("different line should not be duplicate")
	}
	// Same as first line again — not consecutive duplicate now
	if d.IsDuplicate("line") {
		t.Fatal("non-consecutive repeat should not be duplicate in Consecutive mode")
	}
}

func TestGlobal_RemovesAllSeen(t *testing.T) {
	d := New(Global)
	if d.IsDuplicate("hello") {
		t.Fatal("first occurrence should not be duplicate")
	}
	if d.IsDuplicate("world") {
		t.Fatal("first occurrence of second line should not be duplicate")
	}
	if !d.IsDuplicate("hello") {
		t.Fatal("repeated line should be duplicate in Global mode")
	}
}

func TestGlobal_UniqueCount(t *testing.T) {
	d := New(Global)
	input := []string{"a", "b", "a", "c", "b"}
	for _, l := range input {
		d.IsDuplicate(l)
	}
	if got := d.UniqueCount(); got != 3 {
		t.Errorf("UniqueCount = %d, want 3", got)
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := New(Global)
	d.IsDuplicate("line")
	d.Reset()
	if d.IsDuplicate("line") {
		t.Fatal("after Reset, line should not be considered duplicate")
	}
	if d.UniqueCount() != 1 {
		t.Errorf("UniqueCount after Reset = %d, want 1", d.UniqueCount())
	}
}

func TestConsecutive_Reset(t *testing.T) {
	d := New(Consecutive)
	d.IsDuplicate("x")
	d.Reset()
	if d.IsDuplicate("x") {
		t.Fatal("after Reset, line should not be consecutive duplicate")
	}
}
