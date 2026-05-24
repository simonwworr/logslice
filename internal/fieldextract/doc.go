// Package fieldextract provides structured field extraction from log lines.
//
// It supports two common log formats:
//
//   - JSON objects: fields are extracted by key name from the parsed object.
//   - Key=value pairs: fields are extracted by scanning space-delimited tokens
//     of the form key=value or key="quoted value".
//
// JSON parsing is attempted first; key=value scanning is used as a fallback.
//
// Usage:
//
//	e := fieldextract.New(
//		fieldextract.WithKeys("level", "msg", "host"),
//	)
//	fields := e.Extract(line)
//	fmt.Println(fields["level"])
package fieldextract
