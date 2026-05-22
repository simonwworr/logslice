package reader

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"testing"
)

func readerFromString(s string) *LineReader {
	return NewReader(strings.NewReader(s), nil)
}

func TestReadLine_AllLines(t *testing.T) {
	input := "line one\nline two\nline three\n"
	lr := readerFromString(input)
	defer lr.Close()

	expected := []string{"line one", "line two", "line three"}
	for i, want := range expected {
		got, ok := lr.ReadLine()
		if !ok {
			t.Fatalf("expected line %d, got EOF", i)
		}
		if got != want {
			t.Errorf("line %d: got %q, want %q", i, got, want)
		}
	}
	_, ok := lr.ReadLine()
	if ok {
		t.Error("expected EOF after last line")
	}
	if err := lr.Err(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestReadLine_EmptyInput(t *testing.T) {
	lr := readerFromString("")
	defer lr.Close()
	_, ok := lr.ReadLine()
	if ok {
		t.Error("expected no lines from empty input")
	}
}

func TestReadLine_NoTrailingNewline(t *testing.T) {
	lr := readerFromString("hello")
	defer lr.Close()
	line, ok := lr.ReadLine()
	if !ok || line != "hello" {
		t.Errorf("got (%q, %v), want (\"hello\", true)", line, ok)
	}
}

func TestNewReader_GzipTransparent(t *testing.T) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, _ = io.WriteString(gw, "compressed line\n")
	_ = gw.Close()

	gr, err := gzip.NewReader(&buf)
	if err != nil {
		t.Fatal(err)
	}
	lr := NewReader(gr, gr)
	defer lr.Close()

	line, ok := lr.ReadLine()
	if !ok || line != "compressed line" {
		t.Errorf("got (%q, %v), want (\"compressed line\", true)", line, ok)
	}
}
