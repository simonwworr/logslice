package parser

import (
	"testing"
	"time"
)

func TestParseTimestamp_RFC3339(t *testing.T) {
	line := "2024-03-15T10:22:05Z INFO server started"
	ts, offset, err := ParseTimestamp(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 3, 15, 10, 22, 5, 0, time.UTC)
	if !ts.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, ts)
	}
	if offset == 0 {
		t.Error("expected non-zero offset")
	}
}

func TestParseTimestamp_SpaceSeparated(t *testing.T) {
	line := "2024-03-15 10:22:05 ERROR disk full"
	ts, _, err := ParseTimestamp(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts.Year() != 2024 || ts.Month() != 3 || ts.Day() != 15 {
		t.Errorf("unexpected date components in %v", ts)
	}
}

func TestParseTimestamp_NoTimestamp(t *testing.T) {
	line := "this line has no timestamp at all"
	_, _, err := ParseTimestamp(line)
	if err == nil {
		t.Fatal("expected error for line with no timestamp")
	}
	if _, ok := err.(*ErrNoTimestamp); !ok {
		t.Errorf("expected ErrNoTimestamp, got %T", err)
	}
}

func TestParseTimestamp_EmptyLine(t *testing.T) {
	_, _, err := ParseTimestamp("")
	if err == nil {
		t.Fatal("expected error for empty line")
	}
}

func TestDetectFormat(t *testing.T) {
	lines := []string{
		"2024-03-15T10:22:05Z INFO a",
		"2024-03-15T10:22:06Z WARN b",
		"2024-03-15T10:22:07Z ERROR c",
	}
	format := DetectFormat(lines)
	if format == "" {
		t.Error("expected a format to be detected")
	}
}

func TestDetectFormat_NoMatch(t *testing.T) {
	lines := []string{
		"no timestamp here",
		"also no timestamp",
	}
	format := DetectFormat(lines)
	if format != "" {
		t.Errorf("expected empty format, got %q", format)
	}
}
