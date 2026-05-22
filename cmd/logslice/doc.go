// Package main provides the logslice command-line interface.
//
// Usage:
//
//	logslice --from <RFC3339> --to <RFC3339> [--input <file>] [--output <file>]
//
// Flags:
//
//	--from   Start of the time range to extract (inclusive, RFC3339).
//	--to     End of the time range to extract (exclusive, RFC3339).
//	--input  Path to the input log file. Gzip-compressed files (.gz) are
//	         decompressed transparently. Defaults to stdin.
//	--output Path to the output file. Defaults to stdout.
//
// logslice reads the input line by line, parses the leading timestamp on each
// line, and writes only those lines whose timestamp falls within [from, to) to
// the output. Lines without a recognisable timestamp are skipped silently.
//
// Example:
//
//	logslice --from 2024-01-15T10:00:00Z --to 2024-01-15T11:00:00Z \
//	         --input /var/log/app.log.gz --output /tmp/hour.log
package main
