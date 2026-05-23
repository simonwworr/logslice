// Package tail implements a memory-efficient circular buffer for retaining
// the last N lines read from an io.Reader.
//
// It is designed for use in logslice when a user wants to preview the tail
// of a large log file without loading the entire file into memory.
//
// Basic usage:
//
//	tl := tail.New(20)
//	if err := tl.ReadFrom(f); err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range tl.Lines() {
//		fmt.Println(line)
//	}
package tail
