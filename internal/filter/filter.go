// Package filter provides keyword and pattern-based filtering for log lines.
package filter

import (
	"regexp"
	"strings"
)

// Filter holds compiled filtering criteria applied to log lines.
type Filter struct {
	keywords []string
	patterns []*regexp.Regexp
	invert   bool
}

// Option configures a Filter.
type Option func(*Filter)

// WithKeywords adds case-insensitive substring keywords to match against.
func WithKeywords(kws ...string) Option {
	return func(f *Filter) {
		for _, kw := range kws {
			f.keywords = append(f.keywords, strings.ToLower(kw))
		}
	}
}

// WithPatterns adds compiled regular expressions to match against.
func WithPatterns(patterns ...string) (Option, error) {
	var compiled []*regexp.Regexp
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return func(f *Filter) {
		f.patterns = append(f.patterns, compiled...)
	}, nil
}

// WithInvert inverts the filter so only non-matching lines are kept.
func WithInvert(invert bool) Option {
	return func(f *Filter) {
		f.invert = invert
	}
}

// New creates a Filter with the given options.
func New(opts ...Option) *Filter {
	f := &Filter{}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Match reports whether line satisfies the filter criteria.
// A line matches if it satisfies ANY keyword or pattern.
// If no criteria are set, every line matches.
func (f *Filter) Match(line string) bool {
	if len(f.keywords) == 0 && len(f.patterns) == 0 {
		return !f.invert
	}

	lower := strings.ToLower(line)
	for _, kw := range f.keywords {
		if strings.Contains(lower, kw) {
			return !f.invert
		}
	}
	for _, re := range f.patterns {
		if re.MatchString(line) {
			return !f.invert
		}
	}
	return f.invert
}
