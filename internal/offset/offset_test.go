package offset

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func buildIndex(pairs [][2]int64) *Index {
	idx := New()
	for _, p := range pairs {
		idx.Add(p[0], p[1])
	}
	return idx
}

func TestNew_EmptyIndex(t *testing.T) {
	idx := New()
	if idx.Len() != 0 {
		t.Fatalf("expected 0 records, got %d", idx.Len())
	}
}

func TestAdd_IncreasesLen(t *testing.T) {
	idx := New()
	idx.Add(0, 0)
	idx.Add(1000, 4096)
	if idx.Len() != 2 {
		t.Fatalf("expected 2, got %d", idx.Len())
	}
}

func TestLookup_EmptyIndex_ReturnsFalse(t *testing.T) {
	idx := New()
	_, ok := idx.Lookup(100)
	if ok {
		t.Fatal("expected false for empty index")
	}
}

func TestLookup_ExactMatch(t *testing.T) {
	idx := buildIndex([][2]int64{{0, 0}, {1000, 8192}, {2000, 16384}})
	rec, ok := idx.Lookup(1000)
	if !ok {
		t.Fatal("expected record")
	}
	if rec.LineNumber != 1000 || rec.ByteOffset != 8192 {
		t.Fatalf("unexpected record: %+v", rec)
	}
}

func TestLookup_BetweenEntries_ReturnsFloor(t *testing.T) {
	idx := buildIndex([][2]int64{{0, 0}, {1000, 8192}, {2000, 16384}})
	rec, ok := idx.Lookup(1500)
	if !ok {
		t.Fatal("expected record")
	}
	if rec.LineNumber != 1000 {
		t.Fatalf("expected line 1000, got %d", rec.LineNumber)
	}
}

func TestLookup_BeyondLast_ReturnsLast(t *testing.T) {
	idx := buildIndex([][2]int64{{0, 0}, {1000, 8192}})
	rec, ok := idx.Lookup(9999)
	if !ok {
		t.Fatal("expected record")
	}
	if rec.LineNumber != 1000 {
		t.Fatalf("expected line 1000, got %d", rec.LineNumber)
	}
}

func TestRoundTrip_Buffer(t *testing.T) {
	idx := buildIndex([][2]int64{{0, 0}, {500, 2048}, {1000, 4096}})
	var buf bytes.Buffer
	if err := idx.WriteTo(&buf); err != nil {
		t.Fatalf("WriteTo: %v", err)
	}
	got, err := ReadFrom(&buf)
	if err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}
	if got.Len() != idx.Len() {
		t.Fatalf("length mismatch: want %d got %d", idx.Len(), got.Len())
	}
	for i, r := range got.records {
		if r != idx.records[i] {
			t.Fatalf("record %d mismatch: want %+v got %+v", i, idx.records[i], r)
		}
	}
}

func TestSaveAndLoadFile(t *testing.T) {
	idx := buildIndex([][2]int64{{0, 0}, {1000, 8192}})
	path := filepath.Join(t.TempDir(), "test.idx")
	if err := idx.SaveToFile(path); err != nil {
		t.Fatalf("SaveToFile: %v", err)
	}
	loaded, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile: %v", err)
	}
	if loaded.Len() != idx.Len() {
		t.Fatalf("len mismatch")
	}
}

func TestLoadFromFile_NonExistent(t *testing.T) {
	_, err := LoadFromFile(filepath.Join(t.TempDir(), "missing.idx"))
	if !os.IsNotExist(err) {
		// wrapped error — just check it's non-nil
		if err == nil {
			t.Fatal("expected error for missing file")
		}
	}
}
