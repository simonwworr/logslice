// Package columnfilter provides field-level projection for structured log
// lines. It supports two modes:
//
//   - Keep: only the named columns are retained in the output line.
//   - Drop: the named columns are removed; all others pass through.
//
// Both JSON objects and space-separated key=value pairs are handled
// transparently. Lines that do not match either format are returned
// unchanged, making the filter safe to use in mixed-format pipelines.
//
// Example — keep only "level" and "msg" from a JSON log line:
//
//	f := columnfilter.New(
//		columnfilter.WithColumns("level", "msg"),
//		columnfilter.WithMode(columnfilter.Keep),
//	)
//	out, err := f.Apply(line)
package columnfilter
