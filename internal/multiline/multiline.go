// Package multiline provides support for grouping multi-line log entries
// into single logical records before processing.
package multiline

import (
	"regexp"
	"strings"
)

// Joiner accumulates lines that belong to the same logical log entry
// and emits complete records.
type Joiner struct {
	pattern   *regexp.Regexp
	negate    bool
	pending   strings.Builder
	hasPending bool
}

// Option configures a Joiner.
type Option func(*Joiner)

// StartsBlock treats lines matching pattern as the start of a new block.
func StartsBlock(pattern *regexp.Regexp) Option {
	return func(j *Joiner) {
		j.pattern = pattern
		j.negate = false
	}
}

// ContinuesBlock treats lines NOT matching pattern as continuation lines.
func ContinuesBlock(pattern *regexp.Regexp) Option {
	return func(j *Joiner) {
		j.pattern = pattern
		j.negate = true
	}
}

// New creates a Joiner with the provided options.
func New(opts ...Option) *Joiner {
	j := &Joiner{}
	for _, o := range opts {
		o(j)
	}
	return j
}

// Add feeds a raw line to the joiner. If a complete record is ready it is
// returned together with ok=true. When no record is ready yet, ok is false.
// Note: callers must invoke Flush after the last line to retrieve any
// remaining buffered record.
func (j *Joiner) Add(line string) (record string, ok bool) {
	if j.pattern == nil {
		// No grouping configured — pass through immediately.
		return line, true
	}

	matched := j.pattern.MatchString(line)
	isStart := (!j.negate && matched) || (j.negate && !matched)

	if isStart && j.hasPending {
		// Flush previous record before starting a new one.
		record = j.pending.String()
		j.pending.Reset()
		j.pending.WriteString(line)
		return record, true
	}

	if j.hasPending {
		j.pending.WriteByte('\n')
	}
	j.pending.WriteString(line)
	j.hasPending = true
	return "", false
}

// Flush returns any buffered record that has not yet been emitted.
func (j *Joiner) Flush() (record string, ok bool) {
	if !j.hasPending {
		return "", false
	}
	record = j.pending.String()
	j.pending.Reset()
	j.hasPending = false
	return record, true
}

// Reset clears all buffered state.
func (j *Joiner) Reset() {
	j.pending.Reset()
	j.hasPending = false
}

// DrainAll feeds all provided lines through the joiner and returns every
// complete record, including any final buffered record. It is a convenience
// wrapper around repeated calls to Add followed by Flush.
func (j *Joiner) DrainAll(lines []string) []string {
	var records []string
	for _, line := range lines {
		if rec, ok := j.Add(line); ok {
			records = append(records, rec)
		}
	}
	if rec, ok := j.Flush(); ok {
		records = append(records, rec)
	}
	return records
}
