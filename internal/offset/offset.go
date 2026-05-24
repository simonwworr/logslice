// Package offset tracks byte offsets within log files to enable
// fast seeking and resumable processing without re-reading from the start.
package offset

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// Record holds a byte offset and the line number it corresponds to.
type Record struct {
	ByteOffset int64
	LineNumber  int64
}

// Index maps line numbers to byte offsets for a single file.
type Index struct {
	records []Record
}

// New returns an empty Index.
func New() *Index {
	return &Index{}
}

// Add appends a new offset record to the index.
func (idx *Index) Add(lineNumber, byteOffset int64) {
	idx.records = append(idx.records, Record{ByteOffset: byteOffset, LineNumber: lineNumber})
}

// Len returns the number of records in the index.
func (idx *Index) Len() int {
	return len(idx.records)
}

// Lookup returns the Record whose line number is closest to (but not exceeding)
// the requested line. Returns false if the index is empty.
func (idx *Index) Lookup(lineNumber int64) (Record, bool) {
	if len(idx.records) == 0 {
		return Record{}, false
	}
	best := idx.records[0]
	for _, r := range idx.records {
		if r.LineNumber <= lineNumber {
			best = r
		} else {
			break
		}
	}
	return best, true
}

// WriteTo serialises the index into w using a compact binary format.
func (idx *Index) WriteTo(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, int64(len(idx.records))); err != nil {
		return fmt.Errorf("offset: write count: %w", err)
	}
	for _, r := range idx.records {
		if err := binary.Write(w, binary.LittleEndian, r.LineNumber); err != nil {
			return fmt.Errorf("offset: write line: %w", err)
		}
		if err := binary.Write(w, binary.LittleEndian, r.ByteOffset); err != nil {
			return fmt.Errorf("offset: write offset: %w", err)
		}
	}
	return nil
}

// ReadFrom deserialises an index from r.
func ReadFrom(r io.Reader) (*Index, error) {
	var count int64
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return nil, fmt.Errorf("offset: read count: %w", err)
	}
	idx := &Index{records: make([]Record, 0, count)}
	for i := int64(0); i < count; i++ {
		var rec Record
		if err := binary.Read(r, binary.LittleEndian, &rec.LineNumber); err != nil {
			return nil, fmt.Errorf("offset: read line: %w", err)
		}
		if err := binary.Read(r, binary.LittleEndian, &rec.ByteOffset); err != nil {
			return nil, fmt.Errorf("offset: read offset: %w", err)
		}
		idx.records = append(idx.records, rec)
	}
	return idx, nil
}

// SaveToFile writes the index to the given file path, creating or truncating it.
func (idx *Index) SaveToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("offset: create file: %w", err)
	}
	defer f.Close()
	return idx.WriteTo(f)
}

// LoadFromFile reads an index from the given file path.
func LoadFromFile(path string) (*Index, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("offset: open file: %w", err)
	}
	defer f.Close()
	return ReadFrom(f)
}
