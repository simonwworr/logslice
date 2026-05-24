package transform_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/transform"
)

// TestTransformer_PipelineIntegration verifies that a Transformer integrates
// correctly when used as a processing stage: lines are mutated in order and
// the original string is never modified.
func TestTransformer_PipelineIntegration(t *testing.T) {
	original := "  2024-01-15T12:00:00Z ERROR  disk full on /var/log  "

	tr := transform.New(
		transform.WithTrimSpace(),
		transform.WithReplace("ERROR", "WARN"),
		transform.WithStripNonPrintable(),
	)

	result := tr.Apply(original)

	if result == original {
		t.Fatal("expected result to differ from original")
	}

	const want = "2024-01-15T12:00:00Z WARN  disk full on /var/log"
	if result != want {
		t.Fatalf("got %q, want %q", result, want)
	}

	// Original must be unchanged.
	const wantOrig = "  2024-01-15T12:00:00Z ERROR  disk full on /var/log  "
	if original != wantOrig {
		t.Fatalf("original was mutated: %q", original)
	}
}

func TestTransformer_EmptyLine_Stable(t *testing.T) {
	tr := transform.New(
		transform.WithTrimSpace(),
		transform.WithUpperCase(),
	)
	if got := tr.Apply(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}
