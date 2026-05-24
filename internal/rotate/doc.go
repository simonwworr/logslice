// Package rotate discovers and orders rotated log files that share a common
// base path.
//
// Many log rotation schemes append a numeric suffix to the original file name
// when rotating:
//
//	/var/log/app.log      ← current (index 0)
//	/var/log/app.log.1    ← previous rotation (index 1)
//	/var/log/app.log.2    ← older rotation (index 2)
//
// Discover returns all matching files sorted from newest to oldest so that
// callers can open and read them in chronological order (oldest first) or
// reverse order depending on their needs.
//
// Files with non-numeric suffixes (e.g. .gz, .bak) are intentionally ignored
// because they typically require additional decompression handling that is
// outside the scope of this package.
package rotate
