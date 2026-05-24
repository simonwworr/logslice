// Package levelfilter provides log-level based filtering for structured log lines.
// It recognises common severity tokens (DEBUG, INFO, WARN, ERROR, FATAL) and
// discards lines whose level falls below the configured minimum.
package levelfilter

import (
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

// levelNames maps canonical token strings to Level values.
var levelNames = map[string]Level{
	"DEBUG": DEBUG,
	"TRACE": DEBUG, // treat TRACE as DEBUG
	"INFO":  INFO,
	"WARN":  WARN,
	"WARNING": WARN,
	"ERROR": ERROR,
	"ERR":   ERROR,
	"FATAL": FATAL,
	"CRIT":  FATAL,
}

// Filter discards log lines whose detected level is below Min.
type Filter struct {
	Min Level
}

// New returns a Filter that keeps lines at or above min.
func New(min Level) *Filter {
	return &Filter{Min: min}
}

// ParseLevel converts a string such as "warn" or "ERROR" to a Level.
// The second return value is false when the string is not recognised.
func ParseLevel(s string) (Level, bool) {
	l, ok := levelNames[strings.ToUpper(strings.TrimSpace(s))]
	return l, ok
}

// Allow reports whether line should be kept.
// Lines that contain no recognisable level token are always kept so that
// non-structured entries are never silently dropped.
func (f *Filter) Allow(line string) bool {
	upper := strings.ToUpper(line)
	for token, lvl := range levelNames {
		if strings.Contains(upper, token) {
			return lvl >= f.Min
		}
	}
	// No level token found — pass the line through.
	return true
}
