package formatter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/formatter"
)

func TestWriteLine_Raw(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf, formatter.FormatRaw)

	if err := f.WriteLine("hello world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := strings.TrimRight(buf.String(), "\n"); got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestWriteLine_Numbered(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf, formatter.FormatNumbered)

	_ = f.WriteLine("first")
	_ = f.WriteLine("second")

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "1\t") {
		t.Errorf("line 1 should start with '1\\t', got %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "2\t") {
		t.Errorf("line 2 should start with '2\\t', got %q", lines[1])
	}
}

func TestWriteLine_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf, formatter.FormatJSON)

	if err := f.WriteLine("log entry"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"n":1`) {
		t.Errorf("expected JSON with n:1, got %q", out)
	}
	if !strings.Contains(out, `"log entry"`) {
		t.Errorf("expected JSON to contain log entry, got %q", out)
	}
}

func TestLinesWritten(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf, formatter.FormatRaw)

	for i := 0; i < 5; i++ {
		_ = f.WriteLine("line")
	}
	if f.LinesWritten() != 5 {
		t.Errorf("expected 5 lines written, got %d", f.LinesWritten())
	}
}
