// Package formatter provides output formatting for logslice results.
//
// Three formats are supported:
//
//   - FormatRaw      — lines are written verbatim, one per output line.
//   - FormatNumbered — each line is prefixed with its sequential number
//     and a tab character, useful for debugging or piping.
//   - FormatJSON     — each line is wrapped in a minimal JSON envelope
//     containing the line number ("n") and the raw text ("line").
//
// Usage:
//
//	f := formatter.New(os.Stdout, formatter.FormatNumbered)
//	_ = f.WriteLine("2024-01-01T00:00:00Z some event")
//
// The zero value of Format is FormatRaw, so callers that do not need
// a specific format can pass an uninitialized Format value safely.
//
// All Write methods are safe for sequential use but are not
// goroutine-safe; callers that share a Formatter across goroutines
// must synchronise externally.
package formatter
