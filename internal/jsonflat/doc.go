// Package jsonflat provides a Flattener that converts nested JSON log lines
// into a flat map of dot-separated string keys and string values.
//
// This is useful when downstream components (filter, fieldextract, highlight)
// need to operate on deeply nested JSON fields without understanding the
// original structure.
//
// Example:
//
//	f := jsonflat.New(jsonflat.WithSeparator("."))
//	m, err := f.Flatten(`{"http":{"method":"GET","status":200}}`)
//	// m == map[string]string{"http.method": "GET", "http.status": "200"}
//
// By default the separator is ".". Use WithSeparator to change it.
// Use WithMaxDepth to cap recursion; objects beyond the limit are stored as
// their raw JSON representation.
package jsonflat
