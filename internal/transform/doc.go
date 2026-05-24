// Package transform provides a composable line-transformation pipeline for
// log processing.
//
// A Transformer is constructed with zero or more Option functions, each of
// which appends a single string operation to an internal chain. When Apply is
// called the operations are executed in the order they were registered,
// allowing callers to build up arbitrarily complex transformations from small,
// testable primitives.
//
// Built-in options:
//
//	WithTrimSpace        – strip leading/trailing whitespace
//	WithUpperCase        – convert to upper case
//	WithLowerCase        – convert to lower case
//	WithReplace          – find-and-replace all occurrences
//	WithStripNonPrintable – remove non-printable Unicode code points
//
// Custom operations can be added by composing multiple calls to New with
// application-specific option functions that follow the same Option signature.
package transform
