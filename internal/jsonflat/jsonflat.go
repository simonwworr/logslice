// Package jsonflat flattens nested JSON log lines into dot-separated key-value
// pairs, making deeply structured fields accessible for filtering and display.
package jsonflat

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Flattener converts nested JSON objects into a flat map.
type Flattener struct {
	sep       string
	maxDepth  int
}

// Option configures a Flattener.
type Option func(*Flattener)

// WithSeparator sets the key separator (default ".").
func WithSeparator(sep string) Option {
	return func(f *Flattener) { f.sep = sep }
}

// WithMaxDepth limits recursion depth; 0 means unlimited.
func WithMaxDepth(d int) Option {
	return func(f *Flattener) { f.maxDepth = d }
}

// New creates a Flattener with the provided options.
func New(opts ...Option) *Flattener {
	f := &Flattener{sep: ".", maxDepth: 0}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Flatten parses line as JSON and returns a flat map of dot-separated keys to
// string values. Non-JSON lines return an error.
func (f *Flattener) Flatten(line string) (map[string]string, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, fmt.Errorf("jsonflat: not valid JSON: %w", err)
	}
	out := make(map[string]string)
	f.flatten(raw, "", 0, out)
	return out, nil
}

func (f *Flattener) flatten(obj map[string]interface{}, prefix string, depth int, out map[string]string) {
	for k, v := range obj {
		key := k
		if prefix != "" {
			key = prefix + f.sep + k
		}
		switch val := v.(type) {
		case map[string]interface{}:
			if f.maxDepth == 0 || depth+1 < f.maxDepth {
				f.flatten(val, key, depth+1, out)
			} else {
				out[key] = flatJSON(val)
			}
		default:
			out[key] = stringify(v)
		}
	}
}

func stringify(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%v", v), "0"), ".")
	}
}

func flatJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}
