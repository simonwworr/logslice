package truncate

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestNew_NoTruncation_WhenZero(t *testing.T) {
	tr := New(0)
	line := strings.Repeat("a", 200)
	out, cut := tr.Apply(line)
	if cut {
		t.Fatal("expected no truncation for maxBytes=0")
	}
	if out != line {
		t.Fatal("expected original line")
	}
}

func TestApply_ShortLine_Unchanged(t *testing.T) {
	tr := New(100)
	line := "short line"
	out, cut := tr.Apply(line)
	if cut {
		t.Fatalf("expected no cut, got truncated")
	}
	if out != line {
		t.Fatalf("expected %q, got %q", line, out)
	}
}

func TestApply_LongLine_Truncated(t *testing.T) {
	tr := New(20)
	line := strings.Repeat("x", 50)
	out, cut := tr.Apply(line)
	if !cut {
		t.Fatal("expected truncation")
	}
	if len(out) > 20 {
		t.Fatalf("output length %d exceeds maxBytes 20", len(out))
	}
	if !strings.HasSuffix(out, defaultEllipsis) {
		t.Fatalf("expected ellipsis suffix, got %q", out)
	}
}

func TestApply_CustomEllipsis(t *testing.T) {
	tr := New(15, WithEllipsis("[cut]"))
	line := strings.Repeat("a", 40)
	out, cut := tr.Apply(line)
	if !cut {
		t.Fatal("expected truncation")
	}
	if !strings.HasSuffix(out, "[cut]") {
		t.Fatalf("expected [cut] suffix, got %q", out)
	}
	if len(out) > 15 {
		t.Fatalf("output length %d exceeds maxBytes 15", len(out))
	}
}

func TestApply_UTF8_RespectsBoundary(t *testing.T) {
	// Each rune is 3 bytes in UTF-8.
	line := strings.Repeat("あ", 20) // 60 bytes
	tr := New(10)
	out, cut := tr.Apply(line)
	if !cut {
		t.Fatal("expected truncation")
	}
	if !utf8.ValidString(out) {
		t.Fatalf("truncated output is not valid UTF-8: %q", out)
	}
}

func TestTruncated_ReportsCorrectly(t *testing.T) {
	tr := New(10)
	if tr.Truncated("short") {
		t.Fatal("short line should not be flagged as truncated")
	}
	if !tr.Truncated(strings.Repeat("z", 11)) {
		t.Fatal("long line should be flagged as truncated")
	}
}

func TestApply_ExactLength_NotTruncated(t *testing.T) {
	tr := New(10)
	line := strings.Repeat("b", 10)
	_, cut := tr.Apply(line)
	if cut {
		t.Fatal("line exactly at limit should not be truncated")
	}
}
