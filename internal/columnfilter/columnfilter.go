// Package columnfilter selects or drops named columns from structured
// (JSON or key=value) log lines, emitting only the requested fields.
package columnfilter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Mode controls whether the listed columns are kept or dropped.
type Mode int

const (
	Keep Mode = iota // retain only the named columns
	Drop             // remove the named columns, keep the rest
)

// Filter holds the compiled column selection rules.
type Filter struct {
	cols map[string]struct{}
	mode Mode
}

// Option is a functional option for New.
type Option func(*Filter)

// WithColumns specifies the column names to keep or drop.
func WithColumns(names ...string) Option {
	return func(f *Filter) {
		for _, n := range names {
			f.cols[n] = struct{}{}
		}
	}
}

// WithMode sets whether columns are kept or dropped.
func WithMode(m Mode) Option {
	return func(f *Filter) { f.mode = m }
}

// New creates a Filter with the supplied options.
func New(opts ...Option) *Filter {
	f := &Filter{cols: make(map[string]struct{})}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Apply processes a single log line and returns the filtered version.
// It supports JSON objects and key=value pairs. Other formats are
// returned unchanged.
func (f *Filter) Apply(line string) (string, error) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return line, nil
	}
	if line[0] == '{' {
		return f.applyJSON(line)
	}
	return f.applyKV(line), nil
}

func (f *Filter) applyJSON(line string) (string, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, fmt.Errorf("columnfilter: json parse: %w", err)
	}
	out := make(map[string]json.RawMessage, len(obj))
	for k, v := range obj {
		if f.include(k) {
			out[k] = v
		}
	}
	b, err := json.Marshal(out)
	if err != nil {
		return line, fmt.Errorf("columnfilter: json marshal: %w", err)
	}
	return string(b), nil
}

func (f *Filter) applyKV(line string) string {
	parts := strings.Fields(line)
	var kept []string
	for _, p := range parts {
		key := p
		if idx := strings.IndexByte(p, '='); idx >= 0 {
			key = p[:idx]
		}
		if f.include(key) {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, " ")
}

func (f *Filter) include(key string) bool {
	_, found := f.cols[key]
	if f.mode == Keep {
		return found
	}
	return !found
}
