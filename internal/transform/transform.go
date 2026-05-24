package transform

import (
	"strings"
	"unicode"
)

// Option configures a Transformer.
type Option func(*Transformer)

// Transformer applies a chain of string transformations to log lines.
type Transformer struct {
	ops []func(string) string
}

// WithTrimSpace trims leading and trailing whitespace from each line.
func WithTrimSpace() Option {
	return func(t *Transformer) {
		t.ops = append(t.ops, strings.TrimSpace)
	}
}

// WithUpperCase converts each line to upper case.
func WithUpperCase() Option {
	return func(t *Transformer) {
		t.ops = append(t.ops, strings.ToUpper)
	}
}

// WithLowerCase converts each line to lower case.
func WithLowerCase() Option {
	return func(t *Transformer) {
		t.ops = append(t.ops, strings.ToLower)
	}
}

// WithReplace replaces all occurrences of old with new in each line.
func WithReplace(old, new string) Option {
	return func(t *Transformer) {
		t.ops = append(t.ops, func(s string) string {
			return strings.ReplaceAll(s, old, new)
		})
	}
}

// WithStripNonPrintable removes non-printable characters from each line.
func WithStripNonPrintable() Option {
	return func(t *Transformer) {
		t.ops = append(t.ops, func(s string) string {
			return strings.Map(func(r rune) rune {
				if unicode.IsPrint(r) {
					return r
				}
				return -1
			}, s)
		})
	}
}

// New creates a Transformer with the given options.
func New(opts ...Option) *Transformer {
	t := &Transformer{}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply runs all transformation operations on line in order.
// If no operations are configured the original line is returned unchanged.
func (t *Transformer) Apply(line string) string {
	for _, op := range t.ops {
		line = op(line)
	}
	return line
}

// Len returns the number of configured transformation operations.
func (t *Transformer) Len() int {
	return len(t.ops)
}
