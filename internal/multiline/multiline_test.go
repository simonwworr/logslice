package multiline_test

import (
	"regexp"
	"testing"

	"github.com/yourorg/logslice/internal/multiline"
)

func TestNew_NoOptions_PassThrough(t *testing.T) {
	j := multiline.New()
	rec, ok := j.Add("hello")
	if !ok || rec != "hello" {
		t.Fatalf("expected pass-through, got %q ok=%v", rec, ok)
	}
}

func TestAdd_StartsBlock_FlushesOnNewStart(t *testing.T) {
	pat := regexp.MustCompile(`^\d{4}-`)
	j := multiline.New(multiline.StartsBlock(pat))

	_, ok := j.Add("2024-01-01 first line")
	if ok {
		t.Fatal("expected no record on first line")
	}
	_, ok = j.Add("  continuation")
	if ok {
		t.Fatal("expected no record on continuation")
	}
	rec, ok := j.Add("2024-01-02 second entry")
	if !ok {
		t.Fatal("expected flushed record when new block starts")
	}
	if rec != "2024-01-01 first line\n  continuation" {
		t.Errorf("unexpected record: %q", rec)
	}
}

func TestFlush_ReturnsPendingRecord(t *testing.T) {
	pat := regexp.MustCompile(`^\d{4}-`)
	j := multiline.New(multiline.StartsBlock(pat))

	j.Add("2024-03-10 only entry") //nolint:errcheck
	j.Add("  more detail")

	rec, ok := j.Flush()
	if !ok {
		t.Fatal("expected Flush to return pending record")
	}
	if rec != "2024-03-10 only entry\n  more detail" {
		t.Errorf("unexpected flush result: %q", rec)
	}
}

func TestFlush_EmptyJoiner_ReturnsFalse(t *testing.T) {
	j := multiline.New(multiline.StartsBlock(regexp.MustCompile(`^E`)))
	_, ok := j.Flush()
	if ok {
		t.Fatal("expected ok=false on empty flush")
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	pat := regexp.MustCompile(`^START`)
	j := multiline.New(multiline.StartsBlock(pat))
	j.Add("START block")
	j.Reset()
	_, ok := j.Flush()
	if ok {
		t.Fatal("expected empty buffer after Reset")
	}
}

func TestContinuesBlock_GroupsByNonMatch(t *testing.T) {
	// Lines NOT matching the pattern are continuations.
	pat := regexp.MustCompile(`^\d{4}-`)
	j := multiline.New(multiline.ContinuesBlock(pat))

	j.Add("2024-01-01 entry one")
	j.Add("  detail a")
	rec, ok := j.Add("2024-01-02 entry two")
	if !ok {
		t.Fatal("expected record when new timestamped line arrives")
	}
	if rec != "2024-01-01 entry one\n  detail a" {
		t.Errorf("unexpected record: %q", rec)
	}
}
