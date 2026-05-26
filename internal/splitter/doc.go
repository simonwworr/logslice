// Package splitter provides a log-stream splitter that partitions an incoming
// sequence of log lines into multiple numbered output chunks.
//
// Chunks are bounded by either a maximum line count, a maximum byte size, or
// both. When a limit is reached the current chunk is closed and a new one is
// opened via a caller-supplied OpenFunc, giving full control over the
// destination (files, buffers, network writers, etc.).
//
// Typical usage:
//
//	s, err := splitter.New(10_000, 0, func(i int) (io.WriteCloser, error) {
//		return os.Create(fmt.Sprintf("chunk-%04d.log", i))
//	})
//	if err != nil { ... }
//	defer s.Close()
//	for scanner.Scan() {
//		if err := s.WriteLine(scanner.Text()); err != nil { ... }
//	}
package splitter
