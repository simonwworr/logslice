package highlight

import (
	"strings"
	"testing"
)

func TestNew_NoOptions(t *testing.T) {
	h, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Enabled() {
		t.Error("expected Enabled() == false with no rules")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(WithPattern("[invalid", Red))
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	h, _ := New()
	line := "some log line"
	if got := h.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_Keyword_WrapsMatch(t *testing.T) {
	h, err := New(WithKeyword("error", Red))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "an error occurred"
	result := h.Apply(line)
	if !strings.Contains(result, Red+"error"+Reset) {
		t.Errorf("expected highlighted 'error' in %q", result)
	}
	if !strings.Contains(result, "an ") {
		t.Errorf("expected surrounding text preserved in %q", result)
	}
}

func TestApply_Keyword_CaseInsensitive(t *testing.T) {
	h, _ := New(WithKeyword("warn", Yellow))
	result := h.Apply("WARN: disk usage high")
	if !strings.Contains(result, Yellow) {
		t.Errorf("expected yellow highlight for case-insensitive match, got %q", result)
	}
}

func TestApply_Pattern_MultipleMatches(t *testing.T) {
	h, _ := New(WithPattern(`\d+`, Cyan))
	result := h.Apply("line 42 at offset 100")
	// Both numbers should be highlighted.
	count := strings.Count(result, Cyan)
	if count != 2 {
		t.Errorf("expected 2 highlights, got %d in %q", count, result)
	}
}

func TestApply_MultipleRules_EarlierWins(t *testing.T) {
	h, _ := New(
		WithKeyword("error", Red),
		WithKeyword("error", Green),
	)
	result := h.Apply("error")
	// First rule should color it Red; second rule's span is skipped (overlap).
	if !strings.Contains(result, Red) {
		t.Errorf("expected Red from first rule, got %q", result)
	}
	redIdx := strings.Index(result, Red)
	greenIdx := strings.Index(result, Green)
	if greenIdx != -1 && greenIdx < redIdx {
		t.Errorf("Green should not appear before Red")
	}
}

func TestEnabled_AfterAddingRule(t *testing.T) {
	h, _ := New(WithKeyword("info", Cyan))
	if !h.Enabled() {
		t.Error("expected Enabled() == true after adding a rule")
	}
}
