// Package slicer extracts time-range segments from structured log streams.
//
// Given an io.Reader containing log lines with parseable timestamps, Slice
// writes only the lines whose timestamps fall within the specified [From, To]
// window to the provided io.Writer.
//
// Boundary inclusivity is controlled via Options.Inclusive. Lines that cannot
// be parsed for a timestamp are silently skipped, preserving forward
// compatibility with mixed-format logs.
//
// Example usage:
//
//	res, err := slicer.Slice(logFile, os.Stdout, slicer.Options{
//		From:      from,
//		To:        to,
//		Inclusive: true,
//	})
package slicer
