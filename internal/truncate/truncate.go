// Package truncate provides line truncation for long log lines.
// It supports hard truncation at a byte limit and ellipsis insertion
// to signal that content was cut.
package truncate

import "unicode/utf8"

const defaultEllipsis = "..."

// Truncator trims log lines that exceed a configured byte length.
type Truncator struct {
	maxBytes int
	ellipsis string
}

// Option configures a Truncator.
type Option func(*Truncator)

// WithEllipsis sets a custom suffix appended to truncated lines.
func WithEllipsis(e string) Option {
	return func(t *Truncator) {
		t.ellipsis = e
	}
}

// New returns a Truncator that cuts lines to at most maxBytes bytes.
// If maxBytes is <= 0, no truncation is applied.
func New(maxBytes int, opts ...Option) *Truncator {
	t := &Truncator{
		maxBytes: maxBytes,
		ellipsis: defaultEllipsis,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply returns the (possibly truncated) version of line.
// Truncation respects UTF-8 rune boundaries so the result is always valid UTF-8.
// If maxBytes <= 0, the original line is returned unchanged.
func (t *Truncator) Apply(line string) (string, bool) {
	if t.maxBytes <= 0 || len(line) <= t.maxBytes {
		return line, false
	}

	// Reserve space for the ellipsis.
	cutAt := t.maxBytes - len(t.ellipsis)
	if cutAt <= 0 {
		// maxBytes is smaller than the ellipsis itself; just return the ellipsis.
		return t.ellipsis, true
	}

	// Walk back to a valid rune boundary.
	for cutAt > 0 && !utf8.RuneStart(line[cutAt]) {
		cutAt--
	}

	return line[:cutAt] + t.ellipsis, true
}

// Truncated reports whether Apply would truncate the given line.
func (t *Truncator) Truncated(line string) bool {
	return t.maxBytes > 0 && len(line) > t.maxBytes
}
