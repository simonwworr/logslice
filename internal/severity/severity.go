// Package severity maps free-form log level strings to a numeric rank so that
// callers can perform ordered comparisons (e.g. "keep lines >= WARN").
package severity

import (
	"strings"
)

// Level is an ordered numeric severity value.  Higher numbers are more severe.
type Level int

const (
	LevelUnknown Level = iota
	LevelTrace
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelFatal
)

var nameToLevel = map[string]Level{
	"trace":   LevelTrace,
	"debug":   LevelDebug,
	"info":    LevelInfo,
	"notice":  LevelNotice,
	"warn":    LevelWarn,
	"warning": LevelWarn,
	"error":   LevelError,
	"err":     LevelError,
	"fatal":   LevelFatal,
	"crit":    LevelFatal,
	"critical": LevelFatal,
}

var levelToName = map[Level]string{
	LevelUnknown: "UNKNOWN",
	LevelTrace:   "TRACE",
	LevelDebug:   "DEBUG",
	LevelInfo:    "INFO",
	LevelNotice:  "NOTICE",
	LevelWarn:    "WARN",
	LevelError:   "ERROR",
	LevelFatal:   "FATAL",
}

// Parse converts a string such as "WARN" or "warning" to a Level.
// Returns LevelUnknown when the string is not recognised.
func Parse(s string) Level {
	if l, ok := nameToLevel[strings.ToLower(strings.TrimSpace(s))]; ok {
		return l
	}
	return LevelUnknown
}

// String returns the canonical upper-case name for a Level.
func (l Level) String() string {
	if name, ok := levelToName[l]; ok {
		return name
	}
	return "UNKNOWN"
}

// AtLeast reports whether l is at least as severe as min.
func (l Level) AtLeast(min Level) bool {
	return l >= min
}

// Extractor scans a log line and returns the first severity Level found in it.
type Extractor struct{}

// New returns a new Extractor.
func New() *Extractor { return &Extractor{} }

// Extract returns the Level detected in line, or LevelUnknown.
func (e *Extractor) Extract(line string) Level {
	upper := strings.ToUpper(line)
	// Iterate in descending severity so a line containing both "ERROR" and
	// "DEBUG" is classified at the higher level.
	for _, l := range []Level{LevelFatal, LevelError, LevelWarn, LevelNotice, LevelInfo, LevelDebug, LevelTrace} {
		if strings.Contains(upper, levelToName[l]) {
			return l
		}
	}
	return LevelUnknown
}
