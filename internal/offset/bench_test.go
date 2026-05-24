package offset

import (
	"bytes"
	"testing"
)

func BenchmarkAdd_1M(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		idx := New()
		for j := int64(0); j < 1_000_000; j += 1000 {
			idx.Add(j, j*80)
		}
	}
}

func BenchmarkLookup(b *testing.B) {
	idx := New()
	for j := int64(0); j < 1_000_000; j += 1000 {
		idx.Add(j, j*80)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		idx.Lookup(500_000)
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	idx := New()
	for j := int64(0); j < 100_000; j += 1000 {
		idx.Add(j, j*80)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = idx.WriteTo(&buf)
		_, _ = ReadFrom(&buf)
	}
}
