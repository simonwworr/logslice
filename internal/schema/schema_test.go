package schema_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/schema"
)

func TestNew_NoFields_ReturnsError(t *testing.T) {
	_, err := schema.New()
	if err == nil {
		t.Fatal("expected error for empty schema")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := schema.New(schema.WithField("ts", "["))
	if err == nil {
		t.Fatal("expected error for invalid regexp")
	}
}

func TestExtract_MatchesLine(t *testing.T) {
	s, err := schema.New(
		schema.WithField("ts", `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`),
		schema.WithField("level", `[A-Z]+`),
		schema.WithField("msg", `.+`),
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	fields, ok := s.Extract("2024-01-15T08:30:00 INFO user logged in")
	if !ok {
		t.Fatal("expected match")
	}
	if fields["ts"] != "2024-01-15T08:30:00" {
		t.Errorf("ts = %q", fields["ts"])
	}
	if fields["level"] != "INFO" {
		t.Errorf("level = %q", fields["level"])
	}
	if fields["msg"] != "user logged in" {
		t.Errorf("msg = %q", fields["msg"])
	}
}

func TestExtract_NoMatch_ReturnsFalse(t *testing.T) {
	s, _ := schema.New(schema.WithField("num", `\d+`))
	_, ok := s.Extract("no digits here at all")
	if ok {
		t.Fatal("expected no match")
	}
}

func TestExtract_CustomSeparator(t *testing.T) {
	s, err := schema.New(
		schema.WithField("a", `\w+`),
		schema.WithField("b", `\w+`),
		schema.WithSeparator("|"),
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	fields, ok := s.Extract("foo|bar")
	if !ok {
		t.Fatal("expected match")
	}
	if fields["a"] != "foo" || fields["b"] != "bar" {
		t.Errorf("got a=%q b=%q", fields["a"], fields["b"])
	}
}

func TestFields_ReturnsOrderedNames(t *testing.T) {
	s, _ := schema.New(
		schema.WithField("first", `\w+`),
		schema.WithField("second", `\w+`),
	)
	f := s.Fields()
	if len(f) != 2 || f[0] != "first" || f[1] != "second" {
		t.Errorf("unexpected fields: %v", f)
	}
}
