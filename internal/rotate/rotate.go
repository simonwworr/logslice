// Package rotate provides support for reading across rotated log files.
// It discovers and orders log files that share a common base name with
// optional numeric or date-based suffixes (e.g. app.log, app.log.1, app.log.2).
package rotate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// File represents a single rotated log file with its resolved path and
// rotation index (0 = current, higher = older).
type File struct {
	Path  string
	Index int
}

// Discover returns an ordered slice of rotated log files derived from
// basePath. The current file (index 0) comes first; older rotations follow
// in ascending index order. Only files that exist on disk are returned.
func Discover(basePath string) ([]File, error) {
	dir := filepath.Dir(basePath)
	base := filepath.Base(basePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("rotate: read dir %q: %w", dir, err)
	}

	var files []File
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == base {
			files = append(files, File{Path: filepath.Join(dir, name), Index: 0})
			continue
		}
		if strings.HasPrefix(name, base+".") {
			suffix := name[len(base)+1:]
			idx := parseIndex(suffix)
			if idx > 0 {
				files = append(files, File{Path: filepath.Join(dir, name), Index: idx})
			}
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Index < files[j].Index
	})
	return files, nil
}

// parseIndex converts a numeric suffix string to a positive integer.
// Returns 0 if the suffix is not a plain positive integer.
func parseIndex(s string) int {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	if n == 0 {
		return 0
	}
	return n
}
