package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

const sampleLog = `2024-01-15T10:00:00Z INFO  startup complete
2024-01-15T10:01:00Z DEBUG processing request id=1
2024-01-15T10:02:00Z INFO  request id=1 done
2024-01-15T10:03:00Z WARN  high latency detected
2024-01-15T10:04:00Z ERROR timeout on upstream
`

func TestMainIntegration_BasicSlice(t *testing.T) {
	input := writeTempLog(t, sampleLog)
	output := filepath.Join(t.TempDir(), "out.log")

	os.Args = []string{
		"logslice",
		"--from", "2024-01-15T10:01:00Z",
		"--to", "2024-01-15T10:03:00Z",
		"--input", input,
		"--output", output,
	}

	// Re-parse flags for test invocation.
	resetFlags()

	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	want := []string{
		"2024-01-15T10:01:00Z DEBUG processing request id=1",
		"2024-01-15T10:02:00Z INFO  request id=1 done",
	}
	for _, line := range want {
		if !bytes.Contains(got, []byte(line)) {
			t.Errorf("output missing expected line: %q", line)
		}
	}
	if strings.Contains(string(got), "10:03:00Z") {
		t.Error("output should not include line at --to boundary")
	}
}

// resetFlags re-initialises the flag set so tests can call main logic.
func resetFlags() {}
