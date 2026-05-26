package schema_test

import (
	"fmt"
	"testing"

	"github.com/yourorg/logslice/internal/schema"
)

var benchSchema *schema.Schema

func init() {
	var err error
	benchSchema, err = schema.New(
		schema.WithField("ts", `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`),
		schema.WithField("level", `[A-Z]+`),
		schema.WithField("msg", `.+`),
	)
	if err != nil {
		panic(err)
	}
}

func BenchmarkExtract_Match(b *testing.B) {
	line := "2024-06-01T12:00:00 ERROR something went wrong"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = benchSchema.Extract(line)
	}
}

func BenchmarkExtract_NoMatch(b *testing.B) {
	line := "this line has no structure at all"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = benchSchema.Extract(line)
	}
}

func BenchmarkExtract_Parallel(b *testing.B) {
	lines := make([]string, 100)
	for i := range lines {
		lines[i] = fmt.Sprintf("2024-06-01T12:00:%02d INFO msg number %d", i%60, i)
	}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = benchSchema.Extract(lines[i%len(lines)])
			i++
		}
	})
}
