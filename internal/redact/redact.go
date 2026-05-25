// Package redact provides line-level redaction of sensitive patterns
// such as passwords, tokens, and PII before log lines are emitted.
package redact

import (
	"regexp"
	"strings"
)

// Rule describes a single redaction rule: a compiled pattern and the
// replacement string to substitute for every match.
type Rule struct {
	pattern     *regexp.Regexp
	replacement string
}

// Redactor applies a set of redaction rules to log lines.
type Redactor struct {
	rules []Rule
}

// Option is a functional option for configuring a Redactor.
type Option func(*Redactor) error

// WithPattern adds a regex-based redaction rule. Every sub-string
// matching expr is replaced with replacement.
func WithPattern(expr, replacement string) Option {
	return func(r *Redactor) error {
		re, err := regexp.Compile(expr)
		if err != nil {
			return err
		}
		r.rules = append(r.rules, Rule{pattern: re, replacement: replacement})
		return nil
	}
}

// WithKeyword adds a literal-string redaction rule. Matching is
// case-insensitive; the original casing is replaced by replacement.
func WithKeyword(keyword, replacement string) Option {
	return func(r *Redactor) error {
		expr := "(?i)" + regexp.QuoteMeta(keyword)
		re, err := regexp.Compile(expr)
		if err != nil {
			return err
		}
		r.rules = append(r.rules, Rule{pattern: re, replacement: replacement})
		return nil
	}
}

// New constructs a Redactor from the supplied options.
// Returns an error if any option fails (e.g. invalid regex).
func New(opts ...Option) (*Redactor, error) {
	r := &Redactor{}
	for _, o := range opts {
		if err := o(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// Apply runs all redaction rules against line and returns the
// sanitised result. If no rules are configured the original line is
// returned unchanged.
func (r *Redactor) Apply(line string) string {
	if len(r.rules) == 0 {
		return line
	}
	out := line
	for _, rule := range r.rules {
		out = rule.pattern.ReplaceAllString(out, rule.replacement)
	}
	return out
}

// HasRules reports whether any redaction rules have been registered.
func (r *Redactor) HasRules() bool {
	return len(r.rules) > 0
}

// RuleCount returns the number of registered rules.
func (r *Redactor) RuleCount() int {
	return len(r.rules)
}

// MustNew is like New but panics on error. Intended for use in tests
// or program initialisation where the patterns are known-good literals.
func MustNew(opts ...Option) *Redactor {
	r, err := New(opts...)
	if err != nil {
		panic("redact: " + err.Error())
	}
	return r
}

// ensure strings is used (imported for future use in helpers)
var _ = strings.Contains
