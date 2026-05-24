package window

import (
	"testing"
	"time"
)

func ts(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestNew_InvalidSize(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero duration, got nil")
	}
	_, err = New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative duration, got nil")
	}
}

func TestNew_ValidSize(t *testing.T) {
	w, err := New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(w.Buckets()) != 0 {
		t.Fatalf("expected 0 buckets, got %d", len(w.Buckets()))
	}
}

func TestAdd_SingleBucket(t *testing.T) {
	w, _ := New(time.Minute)
	w.Add(ts("2024-01-01T10:00:00Z"), "line1")
	w.Add(ts("2024-01-01T10:00:30Z"), "line2")

	buckets := w.Buckets()
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if buckets[0].Count != 2 {
		t.Errorf("expected count 2, got %d", buckets[0].Count)
	}
}

func TestAdd_MultipleBuckets(t *testing.T) {
	w, _ := New(time.Minute)
	w.Add(ts("2024-01-01T10:00:00Z"), "line1")
	w.Add(ts("2024-01-01T10:01:00Z"), "line2")
	w.Add(ts("2024-01-01T10:02:05Z"), "line3")

	buckets := w.Buckets()
	if len(buckets) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(buckets))
	}
	if buckets[0].Count != 1 || buckets[1].Count != 1 || buckets[2].Count != 1 {
		t.Errorf("unexpected bucket counts: %d %d %d",
			buckets[0].Count, buckets[1].Count, buckets[2].Count)
	}
}

func TestBucket_String(t *testing.T) {
	w, _ := New(time.Minute)
	w.Add(ts("2024-01-01T10:00:00Z"), "hello")
	b := w.Buckets()[0]
	s := b.String()
	if s == "" {
		t.Error("expected non-empty bucket string")
	}
}

func TestReset_ClearsBuckets(t *testing.T) {
	w, _ := New(time.Minute)
	w.Add(ts("2024-01-01T10:00:00Z"), "line1")
	w.Reset()
	if len(w.Buckets()) != 0 {
		t.Fatalf("expected 0 buckets after reset, got %d", len(w.Buckets()))
	}
	// Ensure we can add again after reset
	w.Add(ts("2024-01-01T11:00:00Z"), "line2")
	if len(w.Buckets()) != 1 {
		t.Fatalf("expected 1 bucket after re-add, got %d", len(w.Buckets()))
	}
}
