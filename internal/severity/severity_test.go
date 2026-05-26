package severity_test

import (
	"testing"

	"github.com/logslice/logslice/internal/severity"
)

func TestParse_KnownLevels(t *testing.T) {
	cases := []struct {
		input string
		want  severity.Level
	}{
		{"trace", severity.LevelTrace},
		{"DEBUG", severity.LevelDebug},
		{"Info", severity.LevelInfo},
		{"notice", severity.LevelNotice},
		{"WARN", severity.LevelWarn},
		{"warning", severity.LevelWarn},
		{"error", severity.LevelError},
		{"ERR", severity.LevelError},
		{"fatal", severity.LevelFatal},
		{"CRITICAL", severity.LevelFatal},
		{"crit", severity.LevelFatal},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := severity.Parse(tc.input)
			if got != tc.want {
				t.Errorf("Parse(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParse_Unknown(t *testing.T) {
	for _, s := range []string{"", "verbose", "42", "  "} {
		if got := severity.Parse(s); got != severity.LevelUnknown {
			t.Errorf("Parse(%q) = %v, want LevelUnknown", s, got)
		}
	}
}

func TestLevel_String(t *testing.T) {
	if got := severity.LevelError.String(); got != "ERROR" {
		t.Errorf("String() = %q, want \"ERROR\"", got)
	}
	if got := severity.LevelUnknown.String(); got != "UNKNOWN" {
		t.Errorf("String() = %q, want \"UNKNOWN\"", got)
	}
}

func TestLevel_AtLeast(t *testing.T) {
	if !severity.LevelError.AtLeast(severity.LevelWarn) {
		t.Error("ERROR should be at least WARN")
	}
	if !severity.LevelWarn.AtLeast(severity.LevelWarn) {
		t.Error("WARN should be at least WARN")
	}
	if severity.LevelInfo.AtLeast(severity.LevelWarn) {
		t.Error("INFO should not be at least WARN")
	}
}

func TestExtract_FindsLevelInLine(t *testing.T) {
	e := severity.New()
	cases := []struct {
		line string
		want severity.Level
	}{
		{"2024-01-01 ERROR something went wrong", severity.LevelError},
		{"[WARN] disk usage high", severity.LevelWarn},
		{"level=info msg=started", severity.LevelInfo},
		{"no level here", severity.LevelUnknown},
		{"", severity.LevelUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.line, func(t *testing.T) {
			got := e.Extract(tc.line)
			if got != tc.want {
				t.Errorf("Extract(%q) = %v, want %v", tc.line, got, tc.want)
			}
		})
	}
}

func TestExtract_HigherSeverityWins(t *testing.T) {
	e := severity.New()
	// Line contains both DEBUG and ERROR tokens; ERROR should win.
	line := "DEBUG fallback triggered after ERROR in handler"
	if got := e.Extract(line); got != severity.LevelError {
		t.Errorf("Extract() = %v, want LevelError", got)
	}
}
