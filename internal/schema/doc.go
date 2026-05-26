// Package schema provides a declarative way to define the structure of a log
// line using named field patterns.
//
// A Schema is built from one or more named fields, each backed by a regular
// expression fragment. The fragments are joined with a configurable separator
// (default: space) and compiled into a single anchored regexp. The Extract
// method then parses a raw log line and returns a map of field name to matched
// value.
//
// Example:
//
//	s, err := schema.New(
//		schema.WithField("timestamp", `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`),
//		schema.WithField("level",     `[A-Z]+`),
//		schema.WithField("message",   `.+`),
//	)
//	fields, ok := s.Extract(line)
package schema
