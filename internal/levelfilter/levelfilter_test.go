package levelfilter

import (
	"testing"
)

func TestParseLevel_KnownLevels(t *testing.T) {
	cases := []struct {
		input string
		want  Level
	}{
		{"debug", DEBUG},
		{"INFO", INFO},
		{"Warn", WARN},
		{"WARNING", WARN},
		{"error", ERROR},
		{"ERR", ERROR},
		{"fatal", FATAL},
		{"CRIT", FATAL},
		{"TRACE", DEBUG},
	}
	for _, tc := range cases {
		got, ok := ParseLevel(tc.input)
		if !ok {
			t.Errorf("ParseLevel(%q): expected ok=true", tc.input)
		}
		if got != tc.want {
			t.Errorf("ParseLevel(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

func TestParseLevel_Unknown(t *testing.T) {
	_, ok := ParseLevel("verbose")
	if ok {
		t.Error("expected ok=false for unknown level")
	}
}

func TestAllow_MinInfo_BlocksDebug(t *testing.T) {
	f := New(INFO)
	if f.Allow("2024-01-01T00:00:00Z DEBUG doing something") {
		t.Error("expected DEBUG line to be blocked when min=INFO")
	}
}

func TestAllow_MinInfo_PassesInfo(t *testing.T) {
	f := New(INFO)
	if !f.Allow("2024-01-01T00:00:00Z INFO server started") {
		t.Error("expected INFO line to be allowed when min=INFO")
	}
}

func TestAllow_MinWarn_PassesError(t *testing.T) {
	f := New(WARN)
	if !f.Allow("2024-01-01T00:00:00Z ERROR disk full") {
		t.Error("expected ERROR line to be allowed when min=WARN")
	}
}

func TestAllow_MinError_BlocksWarn(t *testing.T) {
	f := New(ERROR)
	if f.Allow("2024-01-01T00:00:00Z WARN low memory") {
		t.Error("expected WARN line to be blocked when min=ERROR")
	}
}

func TestAllow_NoLevelToken_AlwaysPasses(t *testing.T) {
	f := New(FATAL)
	if !f.Allow("some unstructured log line with no level") {
		t.Error("expected line without level token to always pass")
	}
}

func TestAllow_EmptyLine_Passes(t *testing.T) {
	f := New(ERROR)
	if !f.Allow("") {
		t.Error("expected empty line to pass")
	}
}
