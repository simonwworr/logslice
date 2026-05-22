// Package reader provides line-by-line reading utilities for log files,
// supporting both plain text and gzip-compressed inputs.
package reader

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"strings"
)

// LineReader reads lines from a log source.
type LineReader struct {
	scanner *bufio.Scanner
	closer  io.Closer
}

// NewFileReader opens a log file for reading. Gzip-compressed files
// (identified by a .gz suffix) are decompressed transparently.
func NewFileReader(path string) (*LineReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var r io.Reader = f
	var closer io.Closer = f

	if strings.HasSuffix(path, ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		r = gr
		closer = multiCloser{gr, f}
	}

	return NewReader(r, closer), nil
}

// NewReader wraps an arbitrary io.Reader. closer may be nil.
func NewReader(r io.Reader, closer io.Closer) *LineReader {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	return &LineReader{scanner: scanner, closer: closer}
}

// ReadLine returns the next line and true, or an empty string and false
// when the input is exhausted.
func (lr *LineReader) ReadLine() (string, bool) {
	if lr.scanner.Scan() {
		return lr.scanner.Text(), true
	}
	return "", false
}

// Err returns any scanner error encountered during reading.
func (lr *LineReader) Err() error {
	return lr.scanner.Err()
}

// Close releases any underlying resources.
func (lr *LineReader) Close() error {
	if lr.closer != nil {
		return lr.closer.Close()
	}
	return nil
}

// multiCloser closes multiple io.Closers in order.
type multiCloser []io.Closer

func (mc multiCloser) Close() error {
	var first error
	for _, c := range mc {
		if err := c.Close(); err != nil && first == nil {
			first = err
		}
	}
	return first
}
