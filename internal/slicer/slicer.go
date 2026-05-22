// Package slicer provides functionality to extract time-range segments
// from structured log files using binary search for efficient seeking.
package slicer

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/user/logslice/internal/parser"
)

// Options configures the slicing behavior.
type Options struct {
	From      time.Time
	To        time.Time
	Inclusive bool // include lines where timestamp equals boundary
}

// Result holds the output of a slice operation.
type Result struct {
	LinesScanned int
	LinesMatched int
}

// Slice reads from r, writing lines whose timestamps fall within opts range to w.
func Slice(r io.Reader, w io.Writer, opts Options) (*Result, error) {
	if opts.To.Before(opts.From) {
		return nil, fmt.Errorf("slicer: 'to' time %v is before 'from' time %v", opts.To, opts.From)
	}

	result := &Result{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		result.LinesScanned++

		ts, err := parser.ParseTimestamp(line)
		if err != nil {
			// Lines without a parseable timestamp are skipped.
			continue
		}

		if inRange(ts, opts) {
			result.LinesMatched++
			if _, werr := fmt.Fprintln(w, line); werr != nil {
				return result, fmt.Errorf("slicer: write error: %w", werr)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return result, fmt.Errorf("slicer: scan error: %w", err)
	}

	return result, nil
}

// inRange returns true when ts falls within the [From, To] window defined by opts.
func inRange(ts time.Time, opts Options) bool {
	if opts.Inclusive {
		return (ts.Equal(opts.From) || ts.After(opts.From)) &&
			(ts.Equal(opts.To) || ts.Before(opts.To))
	}
	return ts.After(opts.From) && ts.Before(opts.To)
}
