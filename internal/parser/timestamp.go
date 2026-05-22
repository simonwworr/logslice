package parser

import (
	"fmt"
	"time"
)

// Common log timestamp formats to attempt parsing
var knownFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05.000000",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 2 15:04:05",
	"Jan  2 15:04:05",
}

// ErrNoTimestamp is returned when no timestamp can be parsed from a line.
type ErrNoTimestamp struct {
	Line string
}

func (e *ErrNoTimestamp) Error() string {
	return fmt.Sprintf("no timestamp found in line: %.80s", e.Line)
}

// ParseTimestamp attempts to extract a timestamp from the beginning of a log line.
// It tries each known format and returns the first successful parse along with
// the byte offset where the timestamp ends.
func ParseTimestamp(line string) (time.Time, int, error) {
	for _, format := range knownFormats {
		formatLen := len(format)
		if formatLen > len(line) {
			formatLen = len(line)
		}
		// Try progressively longer substrings to find the timestamp boundary
		for end := formatLen; end <= len(line) && end <= formatLen+5; end++ {
			t, err := time.Parse(format, line[:end])
			if err == nil {
				return t, end, nil
			}
		}
	}
	return time.Time{}, 0, &ErrNoTimestamp{Line: line}
}

// DetectFormat probes a sample of lines to determine the most likely timestamp format.
func DetectFormat(lines []string) string {
	for _, format := range knownFormats {
		hits := 0
		for _, line := range lines {
			formatLen := len(format)
			if formatLen > len(line) {
				continue
			}
			_, err := time.Parse(format, line[:formatLen])
			if err == nil {
				hits++
			}
		}
		if hits >= len(lines)/2 {
			return format
		}
	}
	return ""
}
