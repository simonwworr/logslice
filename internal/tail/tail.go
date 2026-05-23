// Package tail provides functionality for reading the last N lines
// from a log stream, useful for previewing log file endings before slicing.
package tail

import (
	"bufio"
	"container/ring"
	"io"
)

// Tailer holds a fixed-size circular buffer of the most recent lines.
type Tailer struct {
	buf  *ring.Ring
	n    int
	size int
}

// New creates a Tailer that retains the last n lines.
// If n <= 0, New panics.
func New(n int) *Tailer {
	if n <= 0 {
		panic("tail: n must be greater than zero")
	}
	return &Tailer{
		buf: ring.New(n),
		n:   n,
	}
}

// ReadFrom consumes all lines from r, keeping only the last n.
func (t *Tailer) ReadFrom(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		t.buf.Value = line
		t.buf = t.buf.Next()
		if t.size < t.n {
			t.size++
		}
	}
	return scanner.Err()
}

// Lines returns the retained lines in order from oldest to newest.
func (t *Tailer) Lines() []string {
	if t.size == 0 {
		return nil
	}
	result := make([]string, 0, t.size)
	// When buffer is not yet full, unused slots are at the current position
	// and forward. Walk only the filled portion.
	start := t.buf
	if t.size < t.n {
		// Buffer not full: filled entries begin at index 0 (ring start)
		start = t.buf.Move(t.n - t.size).Next().Move(-(t.n - t.size))
	}
	for i := 0; i < t.size; i++ {
		if v := start.Value; v != nil {
			result = append(result, v.(string))
		}
		start = start.Next()
	}
	return result
}

// Reset clears all buffered lines so the Tailer can be reused.
func (t *Tailer) Reset() {
	t.buf = ring.New(t.n)
	t.size = 0
}

// Count returns the number of lines currently buffered.
func (t *Tailer) Count() int {
	return t.size
}
