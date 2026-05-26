// Package splitter divides a log stream into fixed-size or fixed-count chunks,
// each written to a separate output via a user-supplied open function.
package splitter

import (
	"errors"
	"fmt"
	"io"
)

// ErrInvalidChunkSize is returned when the chunk size or line count is not positive.
var ErrInvalidChunkSize = errors.New("splitter: chunk size must be positive")

// OpenFunc is called to obtain a writer for chunk index i.
type OpenFunc func(i int) (io.WriteCloser, error)

// Splitter writes lines to sequentially numbered chunks.
type Splitter struct {
	maxLines int
	maxBytes int64
	open     OpenFunc

	current  io.WriteCloser
	chunkIdx int
	lines    int
	bytes    int64
}

// New creates a Splitter. maxLines and maxBytes are upper bounds per chunk;
either may be zero to disable that limit, but at least one must be positive.
func New(maxLines int, maxBytes int64, open OpenFunc) (*Splitter, error) {
	if maxLines < 0 || maxBytes < 0 {
		return nil, ErrInvalidChunkSize
	}
	if maxLines == 0 && maxBytes == 0 {
		return nil, ErrInvalidChunkSize
	}
	if open == nil {
		return nil, errors.New("splitter: open func must not be nil")
	}
	return &Splitter{maxLines: maxLines, maxBytes: maxBytes, open: open}, nil
}

// WriteLine writes a single log line (without trailing newline) to the current chunk,
// rotating to a new chunk when a limit is reached.
func (s *Splitter) WriteLine(line string) error {
	need := int64(len(line)) + 1 // +1 for newline

	if s.current == nil {
		if err := s.rotate(); err != nil {
			return err
		}
	}

	// Rotate if adding this line would exceed a limit (and chunk is non-empty).
	if s.lines > 0 {
		if s.maxLines > 0 && s.lines >= s.maxLines {
			if err := s.rotate(); err != nil {
				return err
			}
		} else if s.maxBytes > 0 && s.bytes+need > s.maxBytes {
			if err := s.rotate(); err != nil {
				return err
			}
		}
	}

	if _, err := fmt.Fprintln(s.current, line); err != nil {
		return fmt.Errorf("splitter: write chunk %d: %w", s.chunkIdx, err)
	}
	s.lines++
	s.bytes += need
	return nil
}

// Close flushes and closes the current chunk writer.
func (s *Splitter) Close() error {
	if s.current == nil {
		return nil
	}
	err := s.current.Close()
	s.current = nil
	return err
}

// ChunksWritten returns the number of chunks opened so far.
func (s *Splitter) ChunksWritten() int { return s.chunkIdx }

func (s *Splitter) rotate() error {
	if s.current != nil {
		if err := s.current.Close(); err != nil {
			return fmt.Errorf("splitter: close chunk %d: %w", s.chunkIdx, err)
		}
		s.chunkIdx++
	}
	wc, err := s.open(s.chunkIdx)
	if err != nil {
		return fmt.Errorf("splitter: open chunk %d: %w", s.chunkIdx, err)
	}
	s.current = wc
	s.lines = 0
	s.bytes = 0
	return nil
}
