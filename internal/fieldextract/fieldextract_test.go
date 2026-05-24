package fieldextract

import (
	"testing"
)

func TestExtract_JSONLine(t *testing.T) {
	e := New(WithKeys("level", "msg", "ts"))
	got := e.Extract(`{"level":"info","msg":"started","ts":"2024-01-01T00:00:00Z"}`)
	if got["level"] != "info" {
		t.Fatalf("level: want info, got %q", got["level"])
	}
	if got["msg"] != "started" {
		t.Fatalf("msg: want started, got %q", got["msg"])
	}
	if got["ts"] != "2024-01-01T00:00:00Z" {
		t.Fatalf("ts: want 2024-01-01T00:00:00Z, got %q", got["ts"])
	}
}

func TestExtract_KeyValue(t *testing.T) {
	e := New(WithKeys("level", "host"))
	got := e.Extract(`2024-01-01 level=warn host=web-01 msg="disk full"`)
	if got["level"] != "warn" {
		t.Fatalf("level: want warn, got %q", got["level"])
	}
	if got["host"] != "web-01" {
		t.Fatalf("host: want web-01, got %q", got["host"])
	}
}

func TestExtract_MissingKey_NotInResult(t *testing.T) {
	e := New(WithKeys("missing"))
	got := e.Extract(`{"level":"info"}`)
	if _, ok := got["missing"]; ok {
		t.Fatal("expected missing key to be absent from result")
	}
}

func TestExtract_NoKeys_ReturnsEmpty(t *testing.T) {
	e := New()
	got := e.Extract(`{"level":"info","msg":"hello"}`)
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestExtract_NumericJSONValue(t *testing.T) {
	e := New(WithKeys("code"))
	got := e.Extract(`{"code":404,"msg":"not found"}`)
	if got["code"] != "404" {
		t.Fatalf("code: want 404, got %q", got["code"])
	}
}

func TestExtract_QuotedKeyValue(t *testing.T) {
	e := New(WithKeys("env"))
	got := e.Extract(`ts=2024-01-01 env="production" svc=api`)
	if got["env"] != "production" {
		t.Fatalf("env: want production, got %q", got["env"])
	}
}

func TestExtract_EmptyLine(t *testing.T) {
	e := New(WithKeys("level"))
	got := e.Extract("")
	if len(got) != 0 {
		t.Fatalf("expected empty map for empty line, got %v", got)
	}
}
