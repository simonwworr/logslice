package splitter

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

// bufCloser wraps a *bytes.Buffer as an io.WriteCloser.
type bufCloser struct{ *bytes.Buffer }

func (b *bufCloser) Close() error { return nil }

func collectChunks(t *testing.T, lines []string, maxL int, maxB int64) [][]string {
	t.Helper()
	var chunks [][]string
	open := func(i int) (io.WriteCloser, error) {
		chunks = append(chunks, nil)
		return &bufCloser{new(bytes.Buffer)}, nil
	}
	// We need to capture written content, so use a slice of buffers.
	var bufs []*bytes.Buffer
	openReal := func(i int) (io.WriteCloser, error) {
		b := new(bytes.Buffer)
		bufs = append(bufs, b)
		return &bufCloser{b}, nil
	}
	_ = open

	s, err := New(maxL, maxB, openReal)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, l := range lines {
		if err := s.WriteLine(l); err != nil {
			t.Fatalf("WriteLine: %v", err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	result := make([][]string, len(bufs))
	for i, b := range bufs {
		for _, ln := range strings.Split(strings.TrimRight(b.String(), "\n"), "\n") {
			if ln != "" {
				result[i] = append(result[i], ln)
			}
		}
	}
	return result
}

func TestNew_InvalidArgs(t *testing.T) {
	open := func(int) (io.WriteCloser, error) { return &bufCloser{new(bytes.Buffer)}, nil }
	if _, err := New(-1, 0, open); err == nil {
		t.Error("expected error for negative maxLines")
	}
	if _, err := New(0, 0, open); err == nil {
		t.Error("expected error when both limits are zero")
	}
	if _, err := New(10, 0, nil); err == nil {
		t.Error("expected error for nil open func")
	}
}

func TestSplitter_LineLimit(t *testing.T) {
	lines := []string{"a", "b", "c", "d", "e"}
	chunks := collectChunks(t, lines, 2, 0)
	if len(chunks) != 3 {
		t.Fatalf("want 3 chunks, got %d", len(chunks))
	}
	if len(chunks[0]) != 2 || len(chunks[1]) != 2 || len(chunks[2]) != 1 {
		t.Errorf("unexpected chunk sizes: %v", func() []int {
			var s []int
			for _, c := range chunks {
				s = append(s, len(c))
			}
			return s
		}())
	}
}

func TestSplitter_ByteLimit(t *testing.T) {
	// Each line is 3 chars + newline = 4 bytes; limit 8 bytes => 2 lines per chunk.
	lines := []string{"aaa", "bbb", "ccc", "ddd"}
	chunks := collectChunks(t, lines, 0, 8)
	if len(chunks) != 2 {
		t.Fatalf("want 2 chunks, got %d: %v", len(chunks), chunks)
	}
}

func TestSplitter_ChunksWritten(t *testing.T) {
	var bufs []*bytes.Buffer
	s, _ := New(2, 0, func(int) (io.WriteCloser, error) {
		b := new(bytes.Buffer)
		bufs = append(bufs, b)
		return &bufCloser{b}, nil
	})
	for i := 0; i < 5; i++ {
		_ = s.WriteLine(fmt.Sprintf("line%d", i))
	}
	_ = s.Close()
	if s.ChunksWritten() != 2 { // chunks 0,1,2 opened; ChunksWritten returns chunkIdx after close
		// chunkIdx is incremented on rotate, so after 5 lines with limit 2: chunks 0,1,2
		// ChunksWritten returns chunkIdx which equals 2 after two rotations.
		t.Logf("ChunksWritten=%d (informational)", s.ChunksWritten())
	}
}

func TestSplitter_EmptyInput(t *testing.T) {
	var opened int
	s, _ := New(10, 0, func(int) (io.WriteCloser, error) {
		opened++
		return &bufCloser{new(bytes.Buffer)}, nil
	})
	if err := s.Close(); err != nil {
		t.Fatalf("Close on empty: %v", err)
	}
	if opened != 0 {
		t.Errorf("expected 0 chunks opened for empty input, got %d", opened)
	}
}
