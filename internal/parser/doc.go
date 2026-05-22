// Package parser provides utilities for detecting and parsing timestamps
// embedded in structured log lines.
//
// logslice relies on this package to locate time boundaries within large log
// files so that binary search can be used to efficiently seek to the start and
// end of a requested time range without reading the entire file.
//
// Supported timestamp formats include RFC 3339, common syslog formats, and
// several ISO-8601 variants.  Use DetectFormat to probe a sample of lines and
// identify the format in use before processing the full file, or call
// ParseTimestamp directly when the format is unknown and automatic detection
// is preferred.
package parser
