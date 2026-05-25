// Package linenum provides a thread-safe line-number tracker for the
// logslice processing pipeline.
//
// During a slicing run every source line is counted regardless of whether
// it passes filters (absolute counter). Lines that survive all filter and
// slicer stages also increment a separate matched counter.  Both values
// are exposed through Annotation structs so that formatters and output
// writers can prefix lines with either their original file position or
// their position within the result set.
//
// Usage:
//
//	tr := linenum.New()
//	for _, line := range lines {
//		passed := filter.Match(line)
//		ann := tr.Annotate(passed)
//		if passed {
//			fmt.Printf("%d\t%s\n", ann.Relative, line)
//		}
//	}
package linenum
