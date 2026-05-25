package jsonflat

import (
	"testing"
)

func TestFlatten_SimpleObject(t *testing.T) {
	f := New()
	m, err := f.Flatten(`{"level":"info","msg":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["level"] != "info" {
		t.Errorf("level: got %q, want %q", m["level"], "info")
	}
	if m["msg"] != "hello" {
		t.Errorf("msg: got %q, want %q", m["msg"], "hello")
	}
}

func TestFlatten_NestedObject(t *testing.T) {
	f := New()
	m, err := f.Flatten(`{"http":{"method":"GET","status":200}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["http.method"] != "GET" {
		t.Errorf("http.method: got %q, want %q", m["http.method"], "GET")
	}
	if m["http.status"] != "200" {
		t.Errorf("http.status: got %q, want %q", m["http.status"], "200")
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	f := New(WithSeparator("_"))
	m, err := f.Flatten(`{"a":{"b":"c"}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := m["a_b"]; !ok {
		t.Errorf("expected key a_b, got %v", m)
	}
}

func TestFlatten_MaxDepth(t *testing.T) {
	f := New(WithMaxDepth(1))
	m, err := f.Flatten(`{"a":{"b":{"c":"deep"}}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// at depth limit the nested value should be stored as raw JSON
	if _, ok := m["a.b"]; !ok {
		t.Errorf("expected key a.b at depth limit, got %v", m)
	}
}

func TestFlatten_NonJSON_ReturnsError(t *testing.T) {
	f := New()
	_, err := f.Flatten("not json at all")
	if err == nil {
		t.Fatal("expected error for non-JSON input, got nil")
	}
}

func TestFlatten_NullValue(t *testing.T) {
	f := New()
	m, err := f.Flatten(`{"key":null}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["key"] != "" {
		t.Errorf("null should stringify to empty string, got %q", m["key"])
	}
}

func TestFlatten_BoolValue(t *testing.T) {
	f := New()
	m, err := f.Flatten(`{"ok":true,"fail":false}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["ok"] != "true" {
		t.Errorf("ok: got %q, want true", m["ok"])
	}
	if m["fail"] != "false" {
		t.Errorf("fail: got %q, want false", m["fail"])
	}
}
