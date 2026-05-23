package context_test

import (
	"context"
	"testing"
	"time"

	logctx "github.com/example/logslice/internal/context"
)

func TestWithTimeRange_StoresAndRetrieves(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	ctx := logctx.WithTimeRange(context.Background(), start, end)

	gotStart, gotEnd, ok := logctx.TimeRange(ctx)
	if !ok {
		t.Fatal("expected TimeRange to return ok=true")
	}
	if !gotStart.Equal(start) {
		t.Errorf("start: got %v, want %v", gotStart, start)
	}
	if !gotEnd.Equal(end) {
		t.Errorf("end: got %v, want %v", gotEnd, end)
	}
}

func TestTimeRange_NotSet_ReturnsFalse(t *testing.T) {
	_, _, ok := logctx.TimeRange(context.Background())
	if ok {
		t.Error("expected ok=false when time range not set")
	}
}

func TestWithDeadline_CancelsAfterDuration(t *testing.T) {
	ctx, cancel := logctx.WithDeadline(context.Background(), 50*time.Millisecond)
	defer cancel()

	select {
	case <-logctx.Done(ctx):
		if logctx.Err(ctx) == nil {
			t.Error("expected non-nil error after deadline")
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("context was not cancelled within expected time")
	}
}

func TestWithDeadline_CancelFuncStopsContext(t *testing.T) {
	ctx, cancel := logctx.WithDeadline(context.Background(), 10*time.Second)
	cancel()

	select {
	case <-logctx.Done(ctx):
		// expected
	default:
		t.Error("expected context to be done after cancel()")
	}
}

func TestErr_NoError_WhenActive(t *testing.T) {
	ctx := context.Background()
	if err := logctx.Err(ctx); err != nil {
		t.Errorf("expected nil error for active context, got %v", err)
	}
}
