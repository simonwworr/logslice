// Package dedup implements log line deduplication for the logslice pipeline.
//
// Two deduplication modes are supported:
//
//   - Consecutive: suppresses a line only when it immediately follows an
//     identical line. Useful for collapsing repeated error bursts while
//     preserving the same line if it reappears later in the log.
//
//   - Global: suppresses any line that has appeared anywhere earlier in the
//     stream. Suitable for producing a unique-line summary of a log file.
//
// Usage:
//
//	d := dedup.New(dedup.Consecutive)
//	for _, line := range lines {
//	    if !d.IsDuplicate(line) {
//	        fmt.Println(line)
//	    }
//	}
package dedup
