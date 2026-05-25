package columnfilter

import (
	"testing"
)

func TestApply_JSON_KeepColumns(t *testing.T) {
	f := New(WithColumns("level", "msg"), WithMode(Keep))
	out, err := f.Apply(`{"level":"info","msg":"hello","ts":"2024-01-01"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contains(out, "ts") {
		t.Errorf("expected ts to be dropped, got %s", out)
	}
	if !contains(out, "level") || !contains(out, "msg") {
		t.Errorf("expected level and msg to be kept, got %s", out)
	}
}

func TestApply_JSON_DropColumns(t *testing.T) {
	f := New(WithColumns("ts"), WithMode(Drop))
	out, err := f.Apply(`{"level":"info","msg":"hello","ts":"2024-01-01"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contains(out, "ts") {
		t.Errorf("expected ts to be dropped, got %s", out)
	}
	if !contains(out, "level") {
		t.Errorf("expected level to remain, got %s", out)
	}
}

func TestApply_KV_KeepColumns(t *testing.T) {
	f := New(WithColumns("level", "msg"), WithMode(Keep))
	out, err := f.Apply("level=info msg=hello ts=2024-01-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contains(out, "ts") {
		t.Errorf("expected ts dropped, got %q", out)
	}
	if !contains(out, "level=info") || !contains(out, "msg=hello") {
		t.Errorf("expected level and msg kept, got %q", out)
	}
}

func TestApply_KV_DropColumns(t *testing.T) {
	f := New(WithColumns("ts"), WithMode(Drop))
	out, err := f.Apply("level=info msg=hello ts=2024-01-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contains(out, "ts=") {
		t.Errorf("expected ts dropped, got %q", out)
	}
}

func TestApply_EmptyLine_ReturnsEmpty(t *testing.T) {
	f := New(WithColumns("level"), WithMode(Keep))
	out, err := f.Apply("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty, got %q", out)
	}
}

func TestApply_InvalidJSON_ReturnsError(t *testing.T) {
	f := New(WithColumns("level"), WithMode(Keep))
	_, err := f.Apply("{not valid json}")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_NoColumns_KeepMode_DropsAll(t *testing.T) {
	f := New(WithMode(Keep)) // no columns specified
	out, err := f.Apply("level=info msg=hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected all dropped, got %q", out)
	}
}

// contains is a helper to check substring presence.
func contains(s, sub string) bool {
	return len(s) > 0 && len(sub) > 0 &&
		(len(s) >= len(sub)) &&
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})()
}
