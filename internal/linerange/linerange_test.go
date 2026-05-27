package linerange

import (
	"testing"
)

func TestNew_ValidRange(t *testing.T) {
	s, err := New(3, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Start() != 3 || s.End() != 7 {
		t.Fatalf("expected start=3 end=7, got start=%d end=%d", s.Start(), s.End())
	}
}

func TestNew_SingleLine(t *testing.T) {
	s, err := New(5, 5)
	if err != nil {
		t.Fatalf("unexpected error for single-line range: %v", err)
	}
	if !s.InRange(5) {
		t.Fatal("expected line 5 to be in range")
	}
}

func TestNew_InvalidStart_Zero(t *testing.T) {
	_, err := New(0, 5)
	if err == nil {
		t.Fatal("expected error for start=0")
	}
}

func TestNew_InvalidRange_StartGtEnd(t *testing.T) {
	_, err := New(10, 5)
	if err == nil {
		t.Fatal("expected error when start > end")
	}
}

func TestInRange_Boundaries(t *testing.T) {
	s, _ := New(3, 6)
	cases := []struct {
		line    int
		wantIn bool
	}{
		{2, false},
		{3, true},
		{4, true},
		{6, true},
		{7, false},
	}
	for _, c := range cases {
		if got := s.InRange(c.line); got != c.wantIn {
			t.Errorf("InRange(%d) = %v, want %v", c.line, got, c.wantIn)
		}
	}
}

func TestPast_BeyondEnd(t *testing.T) {
	s, _ := New(1, 4)
	if s.Past(4) {
		t.Error("line 4 should not be past end=4")
	}
	if !s.Past(5) {
		t.Error("line 5 should be past end=4")
	}
}

func TestSlice_BasicRange(t *testing.T) {
	lines := []string{"a", "b", "c", "d", "e"}
	s, _ := New(2, 4)
	got := s.Slice(lines)
	want := []string{"b", "c", "d"}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestSlice_RangeExceedsInput(t *testing.T) {
	lines := []string{"x", "y"}
	s, _ := New(1, 100)
	got := s.Slice(lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestSlice_EmptyInput(t *testing.T) {
	s, _ := New(1, 5)
	got := s.Slice([]string{})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d lines", len(got))
	}
}
