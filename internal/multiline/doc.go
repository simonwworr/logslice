// Package multiline groups consecutive log lines that belong to the same
// logical entry into a single record.
//
// Many logging frameworks emit stack traces, JSON payloads, or other
// multi-line content as separate lines after the initial log statement.
// Processing each raw line individually would break timestamp parsing and
// produce incomplete output.
//
// # Usage
//
// Create a Joiner with either StartsBlock (a pattern that marks the first
// line of a new entry) or ContinuesBlock (a pattern that marks lines which
// are NOT the start of a new entry):
//
//	pat := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)
//	j := multiline.New(multiline.StartsBlock(pat))
//
//	for _, raw := range lines {
//		if rec, ok := j.Add(raw); ok {
//			process(rec)
//		}
//	}
//	if rec, ok := j.Flush(); ok {
//		process(rec) // handle the last buffered record
//	}
package multiline
