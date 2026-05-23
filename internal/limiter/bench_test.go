package limiter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/limiter"
)

// BenchmarkAdd_NoLimit measures the overhead of Add when no limits are set.
func BenchmarkAdd_NoLimit(b *testing.B) {
	l, _ := limiter.New(0, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Add(128)
	}
}

// BenchmarkAdd_LineLimit measures Add when a line limit is configured but not
// yet reached (hot path).
func BenchmarkAdd_LineLimit(b *testing.B) {
	l, _ := limiter.New(int64(b.N+1), 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Add(128)
	}
}

// BenchmarkAdd_ByteLimit measures Add when a byte limit is configured but not
// yet reached (hot path).
func BenchmarkAdd_ByteLimit(b *testing.B) {
	l, _ := limiter.New(0, int64(b.N)*128+1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Add(128)
	}
}
