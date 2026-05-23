package highlight

import (
	"regexp"
	"strings"
)

// ANSI color codes.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Cyan   = "\033[36m"
)

// Highlighter applies ANSI color highlighting to log lines based on
// keyword or pattern matches.
type Highlighter struct {
	rules []rule
}

type rule struct {
	pattern *regexp.Regexp
	color   string
}

// Option configures a Highlighter.
type Option func(*Highlighter) error

// WithKeyword adds a case-insensitive literal keyword rule.
func WithKeyword(keyword, color string) Option {
	return func(h *Highlighter) error {
		pat, err := regexp.Compile(`(?i)` + regexp.QuoteMeta(keyword))
		if err != nil {
			return err
		}
		h.rules = append(h.rules, rule{pattern: pat, color: color})
		return nil
	}
}

// WithPattern adds a regex-based highlight rule.
func WithPattern(pattern, color string) Option {
	return func(h *Highlighter) error {
		pat, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		h.rules = append(h.rules, rule{pattern: pat, color: color})
		return nil
	}
}

// New creates a Highlighter with the given options.
// Returns an error if any option fails (e.g. invalid regex).
func New(opts ...Option) (*Highlighter, error) {
	h := &Highlighter{}
	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, err
		}
	}
	return h, nil
}

// Apply returns the line with all matching substrings wrapped in ANSI
// color codes. Rules are applied in order; overlapping matches from
// earlier rules take precedence.
func (h *Highlighter) Apply(line string) string {
	if len(h.rules) == 0 {
		return line
	}

	type span struct {
		start, end int
		color      string
	}

	var spans []span
	for _, r := range h.rules {
		for _, loc := range r.pattern.FindAllStringIndex(line, -1) {
			spans = append(spans, span{loc[0], loc[1], r.color})
		}
	}

	if len(spans) == 0 {
		return line
	}

	var sb strings.Builder
	pos := 0
	for _, s := range spans {
		if s.start < pos {
			continue
		}
		sb.WriteString(line[pos:s.start])
		sb.WriteString(s.color)
		sb.WriteString(line[s.start:s.end])
		sb.WriteString(Reset)
		pos = s.end
	}
	sb.WriteString(line[pos:])
	return sb.String()
}

// Enabled reports whether any highlight rules are configured.
func (h *Highlighter) Enabled() bool {
	return len(h.rules) > 0
}
