package linenum_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/linenum"
)

func BenchmarkAnnotate_Matched(b *testing.B) {
	tr := linenum.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Annotate(true)
	}
}

func BenchmarkAnnotate_Unmatched(b *testing.B) {
	tr := linenum.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Annotate(false)
	}
}

func BenchmarkAnnotate_Parallel(b *testing.B) {
	tr := linenum.New()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tr.Annotate(true)
		}
	})
}
