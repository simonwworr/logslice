package checkpoint_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/checkpoint"
)

func BenchmarkSave(b *testing.B) {
	path := filepath.Join(b.TempDir(), "checkpoint.json")
	s := checkpoint.State{
		File:          "/var/log/app.log",
		Offset:        1 << 20,
		LastTimestamp: time.Now(),
		LinesRead:     50000,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := checkpoint.Save(path, s); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLoad(b *testing.B) {
	path := filepath.Join(b.TempDir(), "checkpoint.json")
	s := checkpoint.State{
		File:          "/var/log/app.log",
		Offset:        1 << 20,
		LastTimestamp: time.Now(),
		LinesRead:     50000,
	}
	if err := checkpoint.Save(path, s); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := checkpoint.Load(path); err != nil {
			b.Fatal(err)
		}
	}
}
