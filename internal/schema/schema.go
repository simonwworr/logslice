package schema

import (
	"fmt"
	"regexp"
	"strings"
)

// Field describes a single named capture group within a log schema.
type Field struct {
	Name    string
	Pattern string
}

// Schema compiles a list of named fields into a single regexp that can
// extract structured key/value pairs from a raw log line.
type Schema struct {
	re     *regexp.Regexp
	fields []string
}

// Option is a functional option for New.
type Option func(*builder)

type builder struct {
	fields []Field
	sep    string
}

// WithField appends a named capture field.
func WithField(name, pattern string) Option {
	return func(b *builder) {
		b.fields = append(b.fields, Field{Name: name, Pattern: pattern})
	}
}

// WithSeparator sets the literal separator inserted between field patterns.
// Defaults to a single space.
func WithSeparator(sep string) Option {
	return func(b *builder) { b.sep = sep }
}

// New constructs a Schema from the supplied options.
// Returns an error if no fields are provided or the compiled regexp is invalid.
func New(opts ...Option) (*Schema, error) {
	b := &builder{sep: " "}
	for _, o := range opts {
		o(b)
	}
	if len(b.fields) == 0 {
		return nil, fmt.Errorf("schema: at least one field is required")
	}

	parts := make([]string, len(b.fields))
	names := make([]string, len(b.fields))
	for i, f := range b.fields {
		parts[i] = fmt.Sprintf("(?P<%s>%s)", f.Name, f.Pattern)
		names[i] = f.Name
	}
	pattern := strings.Join(parts, regexp.QuoteMeta(b.sep))
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("schema: invalid pattern: %w", err)
	}
	return &Schema{re: re, fields: names}, nil
}

// Extract parses line and returns a map of field name → value.
// Returns nil, false when the line does not match the schema.
func (s *Schema) Extract(line string) (map[string]string, bool) {
	m := s.re.FindStringSubmatch(line)
	if m == nil {
		return nil, false
	}
	out := make(map[string]string, len(s.fields))
	for i, name := range s.re.SubexpNames() {
		if name != "" {
			out[name] = m[i]
		}
	}
	return out, true
}

// Fields returns the ordered list of field names defined in the schema.
func (s *Schema) Fields() []string { return append([]string(nil), s.fields...) }
