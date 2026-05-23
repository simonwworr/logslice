// Package highlight provides ANSI terminal color highlighting for log lines.
//
// A Highlighter is configured with one or more rules, each associating a
// keyword or regular expression with an ANSI color code. When Apply is called
// on a log line, every matching substring is wrapped with the corresponding
// color escape sequence and a reset sequence, leaving non-matching text
// unchanged.
//
// Rules are evaluated in declaration order; when two rules would match
// overlapping regions, the earlier rule takes precedence.
//
// Example usage:
//
//	h, err := highlight.New(
//		highlight.WithKeyword("error", highlight.Red),
//		highlight.WithPattern(`\d{4}-\d{2}-\d{2}`, highlight.Cyan),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(h.Apply(line))
package highlight
