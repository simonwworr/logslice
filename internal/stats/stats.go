// Package stats provides collection and reporting of log slicing statistics.
package stats

import (
	"fmt"
	"io"
	"time"
)

// Collector accumulates statistics during a log slicing operation.
type Collector struct {
	LinesRead    int
	LinesMatched int
	LinesFiltered int
	BytesRead    int64
	StartTime    time.Time
	EndTime      time.Time
}

// New returns a new Collector with the start time set to now.
func New() *Collector {
	return &Collector{
		StartTime: time.Now(),
	}
}

// RecordLine records a line that was read from the input.
func (c *Collector) RecordLine(bytes int, matched bool, filtered bool) {
	c.LinesRead++
	c.BytesRead += int64(bytes)
	if matched {
		c.LinesMatched++
	}
	if filtered {
		c.LinesFiltered++
	}
}

// Finish marks the end time of the slicing operation.
func (c *Collector) Finish() {
	c.EndTime = time.Now()
}

// Elapsed returns the duration of the slicing operation.
// If Finish has not been called, it returns the time elapsed so far.
func (c *Collector) Elapsed() time.Duration {
	if c.EndTime.IsZero() {
		return time.Since(c.StartTime)
	}
	return c.EndTime.Sub(c.StartTime)
}

// WriteTo writes a human-readable summary of the statistics to w.
func (c *Collector) WriteTo(w io.Writer) (int64, error) {
	summary := fmt.Sprintf(
		"lines read: %d | matched: %d | filtered: %d | bytes read: %d | elapsed: %s\n",
		c.LinesRead,
		c.LinesMatched,
		c.LinesFiltered,
		c.BytesRead,
		c.Elapsed().Round(time.Millisecond),
	)
	n, err := fmt.Fprint(w, summary)
	return int64(n), err
}
