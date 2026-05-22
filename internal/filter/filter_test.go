package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestMatch_NoFilter_AlwaysTrue(t *testing.T) {
	f := filter.New()
	if !f.Match("any line at all") {
		t.Error("expected match with no criteria")
	}
}

func TestMatch_Keyword_CaseInsensitive(t *testing.T) {
	f := filter.New(filter.WithKeywords("ERROR"))
	if !f.Match("2024-01-01 error: disk full") {
		t.Error("expected keyword match (case-insensitive)")
	}
	if f.Match("2024-01-01 info: all good") {
		t.Error("expected no match for non-matching line")
	}
}

func TestMatch_MultipleKeywords_AnyMatch(t *testing.T) {
	f := filter.New(filter.WithKeywords("warn", "error"))
	if !f.Match("WARNING: low memory") {
		t.Error("expected match on 'warn'")
	}
	if !f.Match("ERROR: timeout") {
		t.Error("expected match on 'error'")
	}
	if f.Match("INFO: startup complete") {
		t.Error("expected no match")
	}
}

func TestMatch_Pattern(t *testing.T) {
	opt, err := filter.WithPatterns(`\bERROR\b`)
	if err != nil {
		t.Fatalf("unexpected compile error: %v", err)
	}
	f := filter.New(opt)
	if !f.Match("ERROR: something failed") {
		t.Error("expected pattern match")
	}
	if f.Match("ERRORS are common") {
		t.Error("expected no match due to word boundary")
	}
}

func TestMatch_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := filter.WithPatterns(`[invalid`)
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestMatch_Invert(t *testing.T) {
	f := filter.New(filter.WithKeywords("debug"), filter.WithInvert(true))
	if f.Match("debug: verbose output") {
		t.Error("expected inverted filter to reject matching line")
	}
	if !f.Match("INFO: normal operation") {
		t.Error("expected inverted filter to accept non-matching line")
	}
}

func TestMatch_Invert_NoFilter(t *testing.T) {
	f := filter.New(filter.WithInvert(true))
	if f.Match("anything") {
		t.Error("expected inverted empty filter to reject everything")
	}
}
