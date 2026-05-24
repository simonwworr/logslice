// Package fieldextract parses structured log fields from a line.
// It supports JSON key extraction and simple key=value pairs.
package fieldextract

import (
	"encoding/json"
	"strings"
)

// Extractor extracts named fields from log lines.
type Extractor struct {
	keys []string
}

// Option configures an Extractor.
type Option func(*Extractor)

// WithKeys specifies which field names to extract.
func WithKeys(keys ...string) Option {
	return func(e *Extractor) {
		e.keys = append(e.keys, keys...)
	}
}

// New creates an Extractor with the given options.
func New(opts ...Option) *Extractor {
	e := &Extractor{}
	for _, o := range opts {
		o(e)
	}
	return e
}

// Extract returns a map of field name → value found in line.
// It attempts JSON parsing first, then falls back to key=value scanning.
func (e *Extractor) Extract(line string) map[string]string {
	result := make(map[string]string, len(e.keys))

	if strings.HasPrefix(strings.TrimSpace(line), "{") {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err == nil {
			for _, k := range e.keys {
				if v, ok := obj[k]; ok {
					result[k] = stringify(v)
				}
			}
			return result
		}
	}

	// key=value fallback
	want := make(map[string]struct{}, len(e.keys))
	for _, k := range e.keys {
		want[k] = struct{}{}
	}
	for _, token := range strings.Fields(line) {
		idx := strings.IndexByte(token, '=')
		if idx <= 0 {
			continue
		}
		k := token[:idx]
		v := strings.Trim(token[idx+1:], `"`)
		if _, ok := want[k]; ok {
			result[k] = v
		}
	}
	return result
}

func stringify(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case nil:
		return ""
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}
