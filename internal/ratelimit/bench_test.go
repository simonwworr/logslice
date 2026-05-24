package ratelimit

import "testing"

func BenchmarkAllow_NoContention(b *testing.B) {
	l, err := New(1e9, 1<<30) // effectively unlimited
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Allow()
	}
}

func BenchmarkAllow_Parallel(b *testing.B) {
	l, err := New(1e9, 1<<30)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Allow()
		}
	})
}
