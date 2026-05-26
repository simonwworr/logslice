// Package severity provides ordered log-level parsing and line-level severity
// extraction for use in logslice pipelines.
//
// # Levels
//
// Levels are defined as ordered constants from LevelTrace (lowest) to
// LevelFatal (highest).  LevelUnknown is returned when no recognised level
// token is present.
//
// # Parsing
//
// Parse converts a free-form string such as "WARN", "warning", or "err" to the
// corresponding Level constant.  Comparison is case-insensitive.
//
// # Extraction
//
// Extractor.Extract scans a raw log line for the first (highest-severity)
// level token it contains, making it suitable for use with structured or
// semi-structured logs that embed the level inline.
package severity
