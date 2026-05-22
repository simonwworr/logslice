// Package formatter provides output formatting utilities for logslice.
package formatter

import (
	"fmt"
	"io"
	"time"
)

// Format represents the output format for sliced log lines.
type Format int

const (
	// FormatRaw outputs lines exactly as they appear in the source.
	FormatRaw Format = iota
	// FormatJSON wraps each line in a JSON envelope with metadata.
	FormatJSON
	// FormatNumbered prefixes each line with its line number.
	FormatNumbered
)

// Formatter writes log lines to an output writer in a specified format.
type Formatter struct {
	w      io.Writer
	fmt    Format
	lineNo int
	start  time.Time
}

// New creates a new Formatter writing to w in the given format.
func New(w io.Writer, f Format) *Formatter {
	return &Formatter{w: w, fmt: f, start: time.Now()}
}

// WriteLine writes a single log line using the configured format.
func (f *Formatter) WriteLine(line string) error {
	f.lineNo++
	switch f.fmt {
	case FormatJSON:
		_, err := fmt.Fprintf(f.w, `{"n":%d,"line":%q}\n`, f.lineNo, line)
		return err
	case FormatNumbered:
		_, err := fmt.Fprintf(f.w, "%d\t%s\n", f.lineNo, line)
		return err
	default:
		_, err := fmt.Fprintln(f.w, line)
		return err
	}
}

// LinesWritten returns the total number of lines written so far.
func (f *Formatter) LinesWritten() int {
	return f.lineNo
}
