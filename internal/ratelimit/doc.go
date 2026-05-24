// Package ratelimit implements a thread-safe token-bucket rate limiter
// intended for use inside the logslice pipeline to cap the number of log
// lines forwarded per second.
//
// # Usage
//
//	l, err := ratelimit.New(1000, 100) // 1 000 lines/s, burst of 100
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range lines {
//		if l.Allow() {
//			fmt.Println(line)
//		}
//	}
//
// Allow is safe for concurrent use.
package ratelimit
