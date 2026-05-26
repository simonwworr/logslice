package headtail

import (
	"fmt"
	"testing"
)

func TestNew_InvalidArgs(t *testing.T) {
	_, err := New(-1, 3)
	if err == nil {
		t.Fatal("expected error for negative headN")
	}
	_, err = New(3, -1)
	if err == nil {
		t.Fatal("expected error for negative tailN")
	}
	_, err = New(0, 0)
	if err == nil {
		t.Fatal("expected error when both are zero")
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := New(3, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLines_HeadOnly(t *testing.T) {
	c, _ := New(3, 0)
	for i := 1; i <= 10; i++ {
		c.Add(fmt.Sprintf("line%d", i))
	}
	got := c.Lines()
	if len(got) != 3 {
		t.Fatalf("want 3 lines, got %d", len(got))
	}
	if got[0] != "line1" || got[2] != "line3" {
		t.Errorf("unexpected head lines: %v", got)
	}
}

func TestLines_TailOnly(t *testing.T) {
	c, _ := New(0, 3)
	for i := 1; i <= 10; i++ {
		c.Add(fmt.Sprintf("line%d", i))
	}
	got := c.Lines()
	if len(got) != 3 {
		t.Fatalf("want 3 lines, got %d", len(got))
	}
	if got[0] != "line8" || got[2] != "line10" {
		t.Errorf("unexpected tail lines: %v", got)
	}
}

func TestLines_HeadAndTail(t *testing.T) {
	c, _ := New(2, 2)
	for i := 1; i <= 8; i++ {
		c.Add(fmt.Sprintf("line%d", i))
	}
	got := c.Lines()
	// expect: line1, line2, line7, line8
	if len(got) != 4 {
		t.Fatalf("want 4, got %d: %v", len(got), got)
	}
	if got[0] != "line1" || got[1] != "line2" || got[2] != "line7" || got[3] != "line8" {
		t.Errorf("unexpected lines: %v", got)
	}
}

func TestLines_FewerThanN(t *testing.T) {
	c, _ := New(5, 5)
	c.Add("a")
	c.Add("b")
	got := c.Lines()
	// head: a,b  tail: a,b  — duplicates preserved
	if len(got) != 4 {
		t.Fatalf("want 4, got %d: %v", len(got), got)
	}
}

func TestTotal(t *testing.T) {
	c, _ := New(2, 2)
	for i := 0; i < 7; i++ {
		c.Add("x")
	}
	if c.Total() != 7 {
		t.Errorf("want 7, got %d", c.Total())
	}
}

func TestReset_ClearsState(t *testing.T) {
	c, _ := New(3, 3)
	for i := 0; i < 6; i++ {
		c.Add("x")
	}
	c.Reset()
	if c.Total() != 0 {
		t.Errorf("expected total 0 after reset, got %d", c.Total())
	}
	if len(c.Lines()) != 0 {
		t.Errorf("expected empty lines after reset")
	}
}
