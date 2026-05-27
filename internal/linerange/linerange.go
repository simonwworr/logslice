// Package linerange provides line-number-based slicing of log content,
// allowing extraction of a contiguous range of lines by index.
package linerange

import "errors"

// ErrInvalidRange is returned when the requested line range is invalid.
var ErrInvalidRange = errors.New("linerange: start must be >= 1 and <= end")

// Slicer extracts lines whose 1-based index falls within [Start, End].
type Slicer struct {
	start int
	end   int
}

// New creates a Slicer that keeps lines in the closed interval [start, end].
// Both start and end are 1-based. Returns ErrInvalidRange if the interval is
// invalid (start < 1 or start > end).
func New(start, end int) (*Slicer, error) {
	if start < 1 || start > end {
		return nil, ErrInvalidRange
	}
	return &Slicer{start: start, end: end}, nil
}

// InRange reports whether the given 1-based line number falls within the
// slicer's interval.
func (s *Slicer) InRange(lineNum int) bool {
	return lineNum >= s.start && lineNum <= s.end
}

// Past reports whether the given 1-based line number is beyond the end of the
// interval, meaning no further lines need to be read.
func (s *Slicer) Past(lineNum int) bool {
	return lineNum > s.end
}

// Slice reads lines from src (a zero-based slice of strings) and returns only
// those whose 1-based position falls within [start, end].
func (s *Slicer) Slice(lines []string) []string {
	out := make([]string, 0)
	for i, line := range lines {
		lineNum := i + 1
		if s.Past(lineNum) {
			break
		}
		if s.InRange(lineNum) {
			out = append(out, line)
		}
	}
	return out
}

// Start returns the first line number of the range (1-based, inclusive).
func (s *Slicer) Start() int { return s.start }

// End returns the last line number of the range (1-based, inclusive).
func (s *Slicer) End() int { return s.end }
