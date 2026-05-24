// Package checkpoint provides lightweight progress persistence for logslice.
//
// When processing large log files it can be useful to resume from where a
// previous run left off rather than re-scanning from the beginning. This
// package serialises a [State] value — containing the source file path, byte
// offset, last seen timestamp and line count — to a JSON file on disk.
//
// Typical usage:
//
//	// Load any existing progress.
//	state, err := checkpoint.Load(".logslice.checkpoint")
//
//	// … process lines, updating state.Offset and state.LinesRead …
//
//	// Persist progress periodically or at the end of the run.
//	checkpoint.Save(".logslice.checkpoint", state)
//
//	// Remove the checkpoint once the run completes successfully.
//	checkpoint.Remove(".logslice.checkpoint")
//
// Saves are atomic on POSIX systems: the state is written to a temporary file
// in the same directory and then renamed into place.
package checkpoint
