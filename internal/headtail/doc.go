// Package headtail implements a head+tail combiner for log streams.
//
// It is useful when you want a quick preview of a large log file without
// reading the entire thing: the first N lines give context about how the
// session started, and the last M lines show the most recent activity.
//
// Usage:
//
//	c, err := headtail.New(10, 10)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for scanner.Scan() {
//		c.Add(scanner.Text())
//	}
//	for _, line := range c.Lines() {
//		fmt.Println(line)
//	}
//
// Lines that fall inside both the head and tail windows (when the total
// number of lines is smaller than headN+tailN) are returned twice so the
// caller always receives a contiguous, ordered slice.
package headtail
