package redact_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/redact"
)

func TestNew_NoOptions_HasNoRules(t *testing.T) {
	r, err := redact.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.HasRules() {
		t.Error("expected no rules")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := redact.New(redact.WithPattern("[", "REDACTED"))
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	r := redact.MustNew()
	line := "user=alice password=secret"
	if got := r.Apply(line); got != line {
		t.Errorf("got %q, want %q", got, line)
	}
}

func TestApply_Pattern_RedactsMatch(t *testing.T) {
	r := redact.MustNew(
		redact.WithPattern(`password=[^\s]+`, "password=[REDACTED]"),
	)
	line := "user=alice password=hunter2 level=info"
	got := r.Apply(line)
	if strings.Contains(got, "hunter2") {
		t.Errorf("sensitive value still present: %q", got)
	}
	if !strings.Contains(got, "password=[REDACTED]") {
		t.Errorf("replacement not found: %q", got)
	}
}

func TestApply_Keyword_CaseInsensitive(t *testing.T) {
	r := redact.MustNew(
		redact.WithKeyword("SECRET", "***"),
	)
	cases := []string{"SECRET", "secret", "Secret", "sEcReT"}
	for _, c := range cases {
		line := "token=" + c
		got := r.Apply(line)
		if strings.Contains(got, c) {
			t.Errorf("keyword %q not redacted in %q", c, got)
		}
	}
}

func TestApply_MultipleRules_AllApplied(t *testing.T) {
	r := redact.MustNew(
		redact.WithPattern(`token=[^\s]+`, "token=[REDACTED]"),
		redact.WithPattern(`\b\d{4}-\d{4}-\d{4}-\d{4}\b`, "[CARD]"),
	)
	line := "token=abc123 card=1234-5678-9012-3456"
	got := r.Apply(line)
	if strings.Contains(got, "abc123") {
		t.Errorf("token not redacted: %q", got)
	}
	if strings.Contains(got, "1234-5678-9012-3456") {
		t.Errorf("card number not redacted: %q", got)
	}
}

func TestRuleCount(t *testing.T) {
	r := redact.MustNew(
		redact.WithKeyword("foo", "***"),
		redact.WithKeyword("bar", "***"),
	)
	if r.RuleCount() != 2 {
		t.Errorf("expected 2 rules, got %d", r.RuleCount())
	}
}

func TestApply_EmptyLine_ReturnsEmpty(t *testing.T) {
	r := redact.MustNew(redact.WithKeyword("secret", "***"))
	if got := r.Apply(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
