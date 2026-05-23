// Package context provides helpers for propagating logslice processing
// state through standard context.Context values.
//
// It supports attaching and retrieving a log time range (start/end) and
// wraps context.WithTimeout for enforcing processing deadlines.
//
// Usage:
//
//	ctx := context.Background()
//	ctx = logctx.WithTimeRange(ctx, start, end)
//
//	start, end, ok := logctx.TimeRange(ctx)
//	if !ok {
//		// time range was not set
//	}
//
//	ctx, cancel := logctx.WithDeadline(ctx, 30*time.Second)
//	defer cancel()
package context
