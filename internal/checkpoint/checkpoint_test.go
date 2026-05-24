package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := tempPath(t)
	now := time.Now().UTC().Truncate(time.Second)
	orig := checkpoint.State{
		File:          "/var/log/app.log",
		Offset:        4096,
		LastTimestamp: now,
		LinesRead:     200,
	}
	if err := checkpoint.Save(path, orig); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := checkpoint.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.File != orig.File {
		t.Errorf("File: got %q, want %q", got.File, orig.File)
	}
	if got.Offset != orig.Offset {
		t.Errorf("Offset: got %d, want %d", got.Offset, orig.Offset)
	}
	if !got.LastTimestamp.Equal(orig.LastTimestamp) {
		t.Errorf("LastTimestamp: got %v, want %v", got.LastTimestamp, orig.LastTimestamp)
	}
	if got.LinesRead != orig.LinesRead {
		t.Errorf("LinesRead: got %d, want %d", got.LinesRead, orig.LinesRead)
	}
	if got.SavedAt.IsZero() {
		t.Error("SavedAt should not be zero after Save")
	}
}

func TestLoad_NonExistentFile_ReturnsZero(t *testing.T) {
	path := tempPath(t)
	state, err := checkpoint.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.File != "" || state.Offset != 0 || state.LinesRead != 0 {
		t.Errorf("expected zero state, got %+v", state)
	}
}

func TestRemove_DeletesFile(t *testing.T) {
	path := tempPath(t)
	if err := checkpoint.Save(path, checkpoint.State{File: "x"}); err != nil {
		t.Fatal(err)
	}
	if err := checkpoint.Remove(path); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should have been deleted")
	}
}

func TestRemove_NonExistent_NoError(t *testing.T) {
	path := tempPath(t)
	if err := checkpoint.Remove(path); err != nil {
		t.Fatalf("Remove on missing file should not error: %v", err)
	}
}
