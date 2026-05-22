// Package reader provides buffered, line-oriented reading of log files for
// the logslice tool.
//
// It supports plain-text files and transparently decompresses gzip-compressed
// files when the path ends with ".gz". The LineReader type wraps a
// bufio.Scanner and exposes a simple ReadLine / Err / Close API that the
// slicer and CLI layers consume without needing to know the underlying
// storage format.
package reader
