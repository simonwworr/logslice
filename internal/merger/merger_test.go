package merger

import (
	"testing"
	"time"
)

func makeSource(entries []Entry) <-chan Entry {
	ch := make(chan Entry, len(entries))
	for _, e := range entries {
		ch <- e
	}
	close(ch)
	return ch
}

func ts(sec int) time.Time {
	return time.Unix(int64(sec), 0).UTC()
}

func TestMerge_SingleSource(t *testing.T) {
	src := makeSource([]Entry{
		{Line: "a", Timestamp: ts(1), Source: 0},
		{Line: "b", Timestamp: ts(3), Source: 0},
		{Line: "c", Timestamp: ts(5), Source: 0},
	})
	m := New(src)
	var got []string
	for e := range m.Merge() {
		got = append(got, e.Line)
	}
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Fatalf("unexpected order: %v", got)
	}
}

func TestMerge_TwoInterleavedSources(t *testing.T) {
	src0 := makeSource([]Entry{
		{Line: "a", Timestamp: ts(1), Source: 0},
		{Line: "c", Timestamp: ts(3), Source: 0},
	})
	src1 := makeSource([]Entry{
		{Line: "b", Timestamp: ts(2), Source: 1},
		{Line: "d", Timestamp: ts(4), Source: 1},
	})
	m := New(src0, src1)
	var got []string
	for e := range m.Merge() {
		got = append(got, e.Line)
	}
	want := []string{"a", "b", "c", "d"}
	for i, w := range want {
		if got[i] != w {
			t.Fatalf("position %d: got %q, want %q", i, got[i], w)
		}
	}
}

func TestMerge_EmptySources(t *testing.T) {
	src := makeSource(nil)
	m := New(src)
	var count int
	for range m.Merge() {
		count++
	}
	if count != 0 {
		t.Fatalf("expected 0 entries, got %d", count)
	}
}

func TestMerge_NoSources(t *testing.T) {
	m := New()
	var count int
	for range m.Merge() {
		count++
	}
	if count != 0 {
		t.Fatalf("expected 0 entries from empty merger, got %d", count)
	}
}
