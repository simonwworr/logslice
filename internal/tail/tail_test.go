package tail

import (
	"strings"
	"testing"
)

func tailerFromString(t *testing.T, n int, input string) *Tailer {
	t.Helper()
	tl := New(n)
	if err := tl.ReadFrom(strings.NewReader(input)); err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}
	return tl
}

func TestNew_InvalidN(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for n=0")
		}
	}()
	New(0)
}

func TestLines_FewerThanN(t *testing.T) {
	input := "line1\nline2\n"
	tl := tailerFromString(t, 5, input)
	got := tl.Lines()
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
	if got[0] != "line1" || got[1] != "line2" {
		t.Errorf("unexpected lines: %v", got)
	}
}

func TestLines_ExactlyN(t *testing.T) {
	input := "a\nb\nc\n"
	tl := tailerFromString(t, 3, input)
	got := tl.Lines()
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
}

func TestLines_MoreThanN(t *testing.T) {
	lines := []string{"one", "two", "three", "four", "five"}
	input := strings.Join(lines, "\n") + "\n"
	tl := tailerFromString(t, 3, input)
	got := tl.Lines()
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(got), got)
	}
	if got[0] != "three" || got[1] != "four" || got[2] != "five" {
		t.Errorf("wrong tail lines: %v", got)
	}
}

func TestLines_EmptyInput(t *testing.T) {
	tl := tailerFromString(t, 5, "")
	if tl.Lines() != nil {
		t.Error("expected nil for empty input")
	}
	if tl.Count() != 0 {
		t.Errorf("expected count 0, got %d", tl.Count())
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	tl := tailerFromString(t, 3, "x\ny\nz\n")
	if tl.Count() != 3 {
		t.Fatalf("expected 3 before reset")
	}
	tl.Reset()
	if tl.Count() != 0 {
		t.Errorf("expected 0 after reset, got %d", tl.Count())
	}
	if tl.Lines() != nil {
		t.Error("expected nil lines after reset")
	}
}

func TestCount_Accuracy(t *testing.T) {
	tl := tailerFromString(t, 10, "a\nb\nc\n")
	if tl.Count() != 3 {
		t.Errorf("expected count 3, got %d", tl.Count())
	}
}
