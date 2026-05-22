package slicer_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/slicer"
)

const sampleLog = `2024-01-15T10:00:00Z INFO  service started
2024-01-15T10:01:00Z DEBUG request received id=1
2024-01-15T10:02:00Z INFO  processed request id=1
2024-01-15T10:03:00Z WARN  slow response id=2 latency=2s
2024-01-15T10:04:00Z ERROR timeout id=3
2024-01-15T10:05:00Z INFO  service stopping
`

func mustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSlice_BasicRange(t *testing.T) {
	opts := slicer.Options{
		From:      mustParse("2024-01-15T10:01:00Z"),
		To:        mustParse("2024-01-15T10:03:00Z"),
		Inclusive: true,
	}

	var out strings.Builder
	res, err := slicer.Slice(strings.NewReader(sampleLog), &out, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LinesMatched != 3 {
		t.Errorf("expected 3 matched lines, got %d", res.LinesMatched)
	}
	if res.LinesScanned != 6 {
		t.Errorf("expected 6 scanned lines, got %d", res.LinesScanned)
	}
}

func TestSlice_ExclusiveBoundaries(t *testing.T) {
	opts := slicer.Options{
		From:      mustParse("2024-01-15T10:01:00Z"),
		To:        mustParse("2024-01-15T10:03:00Z"),
		Inclusive: false,
	}

	var out strings.Builder
	res, err := slicer.Slice(strings.NewReader(sampleLog), &out, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LinesMatched != 1 {
		t.Errorf("expected 1 matched line, got %d", res.LinesMatched)
	}
}

func TestSlice_InvalidRange(t *testing.T) {
	opts := slicer.Options{
		From: mustParse("2024-01-15T10:05:00Z"),
		To:   mustParse("2024-01-15T10:01:00Z"),
	}

	var out strings.Builder
	_, err := slicer.Slice(strings.NewReader(sampleLog), &out, opts)
	if err == nil {
		t.Error("expected error for inverted range, got nil")
	}
}

func TestSlice_EmptyInput(t *testing.T) {
	opts := slicer.Options{
		From:      mustParse("2024-01-15T10:00:00Z"),
		To:        mustParse("2024-01-15T11:00:00Z"),
		Inclusive: true,
	}

	var out strings.Builder
	res, err := slicer.Slice(strings.NewReader(""), &out, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LinesMatched != 0 || res.LinesScanned != 0 {
		t.Errorf("expected zero counts on empty input, got scanned=%d matched=%d",
			res.LinesScanned, res.LinesMatched)
	}
}
