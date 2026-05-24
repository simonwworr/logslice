package fieldextract

import "testing"

var sink map[string]string

func BenchmarkExtract_JSON(b *testing.B) {
	e := New(WithKeys("level", "msg", "ts", "host"))
	line := `{"level":"info","msg":"request handled","ts":"2024-06-01T12:00:00Z","host":"web-01","latency":42}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = e.Extract(line)
	}
}

func BenchmarkExtract_KeyValue(b *testing.B) {
	e := New(WithKeys("level", "msg", "host"))
	line := `2024-06-01T12:00:00Z level=info msg="request handled" host=web-01 latency=42`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = e.Extract(line)
	}
}
