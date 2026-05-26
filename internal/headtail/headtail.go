// Package headtail provides a combiner that emits the first N and last M lines
// from a stream, useful for previewing large log files without reading everything.
package headtail

import "errors"

// Combiner holds head lines as they arrive and delegates tail lines to a
// ring-buffer so the final M lines are always available.
type Combiner struct {
	headN    int
	tailN    int
	headBuf  []string
	tailBuf  []string
	tailPos  int
	tailFull bool
	total    int
}

// New creates a Combiner that keeps the first headN and last tailN lines.
// Either value may be zero to disable that side. Both cannot be zero.
func New(headN, tailN int) (*Combiner, error) {
	if headN < 0 || tailN < 0 {
		return nil, errors.New("headtail: headN and tailN must be non-negative")
	}
	if headN == 0 && tailN == 0 {
		return nil, errors.New("headtail: at least one of headN or tailN must be > 0")
	}
	c := &Combiner{
		headN: headN,
		tailN: tailN,
	}
	if tailN > 0 {
		c.tailBuf = make([]string, tailN)
	}
	return c, nil
}

// Add records a line into the head buffer and/or the tail ring-buffer.
func (c *Combiner) Add(line string) {
	if c.headN > 0 && len(c.headBuf) < c.headN {
		c.headBuf = append(c.headBuf, line)
	}
	if c.tailN > 0 {
		c.tailBuf[c.tailPos] = line
		c.tailPos = (c.tailPos + 1) % c.tailN
		if !c.tailFull && c.tailPos == 0 {
			c.tailFull = true
		}
	}
	c.total++
}

// Lines returns the collected head lines followed by tail lines.
// Duplicate lines that appear in both head and tail are preserved as-is.
func (c *Combiner) Lines() []string {
	out := make([]string, 0, c.headN+c.tailN)
	out = append(out, c.headBuf...)

	if c.tailN > 0 {
		var tail []string
		if c.tailFull {
			tail = append(c.tailBuf[c.tailPos:], c.tailBuf[:c.tailPos]...)
		} else {
			tail = c.tailBuf[:c.tailPos]
		}
		out = append(out, tail...)
	}
	return out
}

// Total returns the number of lines seen via Add.
func (c *Combiner) Total() int { return c.total }

// Reset clears all buffered state so the Combiner can be reused.
func (c *Combiner) Reset() {
	c.headBuf = c.headBuf[:0]
	for i := range c.tailBuf {
		c.tailBuf[i] = ""
	}
	c.tailPos = 0
	c.tailFull = false
	c.total = 0
}
