package transform

import (
	"fmt"
	"testing"
)

func BenchmarkApply_NoOps(b *testing.B) {
	tr := New()
	line := "2024-01-15T12:00:00Z INFO  service started successfully"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tr.Apply(line)
	}
}

func BenchmarkApply_TrimAndLower(b *testing.B) {
	tr := New(WithTrimSpace(), WithLowerCase())
	line := "  2024-01-15T12:00:00Z INFO  service started successfully  "
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tr.Apply(line)
	}
}

func BenchmarkApply_FullChain(b *testing.B) {
	tr := New(
		WithTrimSpace(),
		WithStripNonPrintable(),
		WithReplace("INFO", "DBG"),
		WithLowerCase(),
	)
	lines := make([]string, 1000)
	for i := range lines {
		lines[i] = fmt.Sprintf("  2024-01-15T12:00:00Z INFO request id=%d  ", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tr.Apply(lines[i%len(lines)])
	}
}
