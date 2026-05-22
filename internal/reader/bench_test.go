package reader

import (
	"fmt"
	"strings"
	"testing"
)

// generateLines builds a synthetic log corpus of n lines.
func generateLines(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "2024-01-15T10:%02d:%02dZ INFO request id=%d status=200\n",
			(i/60)%60, i%60, i)
	}
	return sb.String()
}

func BenchmarkReadLine_10k(b *testing.B) {
	corpus := generateLines(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lr := readerFromString(corpus)
		count := 0
		for {
			_, ok := lr.ReadLine()
			if !ok {
				break
			}
			count++
		}
		lr.Close()
		if count != 10_000 {
			b.Fatalf("expected 10000 lines, got %d", count)
		}
	}
}

func BenchmarkReadLine_100k(b *testing.B) {
	corpus := generateLines(100_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lr := readerFromString(corpus)
		for {
			_, ok := lr.ReadLine()
			if !ok {
				break
			}
		}
		lr.Close()
	}
}
