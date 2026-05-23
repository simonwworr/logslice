// Package limiter provides configurable line-count and byte-count limits for
// log slicing operations.
//
// A Limiter is created with optional upper bounds on the number of output lines
// and the total number of output bytes. Callers pass each line's length to Add;
// once either threshold is crossed Add returns ErrLimitReached, signalling that
// the pipeline should stop emitting lines.
//
// Both limits are independent: set a limit to 0 to disable it.
//
// Example:
//
//	lim, err := limiter.New(1000, 0) // stop after 1 000 lines
//	if err != nil { … }
//	for _, line := range lines {
//		if err := lim.Add(int64(len(line))); errors.Is(err, limiter.ErrLimitReached) {
//			break
//		}
//		fmt.Println(line)
//	}
package limiter
