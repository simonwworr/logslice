// Package redact implements configurable redaction of sensitive data
// within log lines before they are written to any output sink.
//
// # Overview
//
// A [Redactor] holds an ordered list of redaction [Rule] values. Each
// rule pairs a compiled regular expression with a replacement string.
// When [Redactor.Apply] is called on a log line, every rule is applied
// in registration order so that later rules can act on already-redacted
// output.
//
// # Building a Redactor
//
// Use [New] with one or more functional options:
//
//	r, err := redact.New(
//		redact.WithPattern(`password=[^\s]+`, "password=[REDACTED]"),
//		redact.WithKeyword("s3cr3t", "***"),
//	)
//
// [WithPattern] accepts any RE2-compatible regular expression.
// [WithKeyword] performs a case-insensitive literal match.
//
// # Thread Safety
//
// A Redactor is safe for concurrent use after construction.
package redact
