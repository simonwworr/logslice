package transform

import (
	"strings"
	"testing"
)

func TestNew_NoOptions_IdentityTransform(t *testing.T) {
	tr := New()
	if tr.Len() != 0 {
		t.Fatalf("expected 0 ops, got %d", tr.Len())
	}
	if got := tr.Apply("hello"); got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}

func TestWithTrimSpace(t *testing.T) {
	tr := New(WithTrimSpace())
	if got := tr.Apply("  hello world  "); got != "hello world" {
		t.Fatalf("unexpected result: %q", got)
	}
}

func TestWithUpperCase(t *testing.T) {
	tr := New(WithUpperCase())
	if got := tr.Apply("hello"); got != "HELLO" {
		t.Fatalf("unexpected result: %q", got)
	}
}

func TestWithLowerCase(t *testing.T) {
	tr := New(WithLowerCase())
	if got := tr.Apply("HELLO"); got != "hello" {
		t.Fatalf("unexpected result: %q", got)
	}
}

func TestWithReplace(t *testing.T) {
	tr := New(WithReplace("foo", "bar"))
	if got := tr.Apply("foo baz foo"); got != "bar baz bar" {
		t.Fatalf("unexpected result: %q", got)
	}
}

func TestWithStripNonPrintable(t *testing.T) {
	tr := New(WithStripNonPrintable())
	input := "hello\x00world\x1b"
	got := tr.Apply(input)
	if strings.ContainsAny(got, "\x00\x1b") {
		t.Fatalf("non-printable chars not stripped: %q", got)
	}
	if got != "helloworld" {
		t.Fatalf("unexpected result: %q", got)
	}
}

func TestChainedOperations(t *testing.T) {
	tr := New(WithTrimSpace(), WithUpperCase(), WithReplace("FOO", "BAR"))
	if got := tr.Apply("  foo baz  "); got != "BAR BAZ" {
		t.Fatalf("unexpected result: %q", got)
	}
}

func TestLen_ReflectsOptionCount(t *testing.T) {
	tr := New(WithTrimSpace(), WithLowerCase())
	if tr.Len() != 2 {
		t.Fatalf("expected 2 ops, got %d", tr.Len())
	}
}
