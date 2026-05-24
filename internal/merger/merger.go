// Package merger provides utilities for merging multiple sorted log streams
// into a single chronologically ordered output.
package merger

import (
	"container/heap"
	"time"
)

// Entry represents a single log line with its parsed timestamp and source index.
type Entry struct {
	Line      string
	Timestamp time.Time
	Source    int
}

// entryHeap implements heap.Interface for min-heap ordering by timestamp.
type entryHeap []Entry

func (h entryHeap) Len() int            { return len(h) }
func (h entryHeap) Less(i, j int) bool  { return h[i].Timestamp.Before(h[j].Timestamp) }
func (h entryHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *entryHeap) Push(x interface{}) { *h = append(*h, x.(Entry)) }
func (h *entryHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Merger merges multiple sorted entry channels into one ordered stream.
type Merger struct {
	sources []<-chan Entry
}

// New creates a Merger that will merge the provided source channels.
// Each source channel must emit entries in ascending timestamp order.
func New(sources ...(<-chan Entry)) *Merger {
	return &Merger{sources: sources}
}

// Merge reads from all source channels and emits entries in timestamp order
// to the returned channel. The returned channel is closed when all sources
// are exhausted.
func (m *Merger) Merge() <-chan Entry {
	out := make(chan Entry, 64)
	go func() {
		defer close(out)
		h := &entryHeap{}
		heap.Init(h)

		// Seed heap with the first entry from each source.
		for _, src := range m.sources {
			if e, ok := <-src; ok {
				heap.Push(h, e)
			}
		}

		for h.Len() > 0 {
			smallest := heap.Pop(h).(Entry)
			out <- smallest
			// Refill from the same source.
			if smallest.Source >= 0 && smallest.Source < len(m.sources) {
				if e, ok := <-m.sources[smallest.Source]; ok {
					heap.Push(h, e)
				}
			}
		}
	}()
	return out
}
