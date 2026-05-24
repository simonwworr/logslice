package rotate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/rotate"
)

func makeFiles(t *testing.T, dir string, names ...string) {
	t.Helper()
	for _, n := range names {
		f, err := os.Create(filepath.Join(dir, n))
		if err != nil {
			t.Fatalf("create %q: %v", n, err)
		}
		f.Close()
	}
}

func TestDiscover_CurrentOnly(t *testing.T) {
	dir := t.TempDir()
	makeFiles(t, dir, "app.log")

	files, err := rotate.Discover(filepath.Join(dir, "app.log"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("want 1 file, got %d", len(files))
	}
	if files[0].Index != 0 {
		t.Errorf("want index 0, got %d", files[0].Index)
	}
}

func TestDiscover_MultipleRotations(t *testing.T) {
	dir := t.TempDir()
	makeFiles(t, dir, "app.log", "app.log.1", "app.log.2", "app.log.3")

	files, err := rotate.Discover(filepath.Join(dir, "app.log"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 4 {
		t.Fatalf("want 4 files, got %d", len(files))
	}
	for i, f := range files {
		if f.Index != i {
			t.Errorf("position %d: want index %d, got %d", i, i, f.Index)
		}
	}
}

func TestDiscover_IgnoresNonNumericSuffix(t *testing.T) {
	dir := t.TempDir()
	makeFiles(t, dir, "app.log", "app.log.gz", "app.log.bak", "app.log.1")

	files, err := rotate.Discover(filepath.Join(dir, "app.log"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// only app.log (index 0) and app.log.1 (index 1) should match
	if len(files) != 2 {
		t.Fatalf("want 2 files, got %d", len(files))
	}
}

func TestDiscover_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	files, err := rotate.Discover(filepath.Join(dir, "app.log"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("want 0 files, got %d", len(files))
	}
}

func TestDiscover_BadDir(t *testing.T) {
	_, err := rotate.Discover("/nonexistent/path/app.log")
	if err == nil {
		t.Fatal("expected error for nonexistent directory")
	}
}
