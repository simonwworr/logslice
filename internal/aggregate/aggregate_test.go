package aggregate

import (
	"testing"
	"time"
)

func ts(base time.Time, offset time.Duration) time.Time {
	return base.Add(offset)
}

func TestNew_InvalidDuration(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
	_, err = New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestNew_ValidDuration(t *testing.T) {
	a, err := New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil aggregator")
	}
}

func TestAdd_ZeroTimestamp_Ignored(t *testing.T) {
	a, _ := New(time.Minute)
	a.Add(time.Time{}, 100)
	if got := a.Buckets(); len(got) != 0 {
		t.Fatalf("expected 0 buckets, got %d", len(got))
	}
}

func TestAdd_SingleBucket(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	a, _ := New(time.Minute)

	a.Add(ts(base, 0), 10)
	a.Add(ts(base, 10*time.Second), 20)
	a.Add(ts(base, 59*time.Second), 5)

	buckets := a.Buckets()
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if buckets[0].Lines != 3 {
		t.Errorf("expected 3 lines, got %d", buckets[0].Lines)
	}
	if buckets[0].Bytes != 35 {
		t.Errorf("expected 35 bytes, got %d", buckets[0].Bytes)
	}
}

func TestAdd_MultipleBuckets(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	a, _ := New(time.Minute)

	a.Add(ts(base, 0), 5)
	a.Add(ts(base, 61*time.Second), 8)
	a.Add(ts(base, 125*time.Second), 3)

	buckets := a.Buckets()
	if len(buckets) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(buckets))
	}
	if buckets[0].Lines != 1 || buckets[0].Bytes != 5 {
		t.Errorf("bucket 0 mismatch: %+v", buckets[0])
	}
	if buckets[1].Lines != 1 || buckets[1].Bytes != 8 {
		t.Errorf("bucket 1 mismatch: %+v", buckets[1])
	}
}

func TestReset_ClearsState(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	a, _ := New(time.Minute)
	a.Add(base, 50)
	a.Reset()
	if got := a.Buckets(); len(got) != 0 {
		t.Fatalf("expected 0 buckets after reset, got %d", len(got))
	}
}

func TestBuckets_ReturnsCopy(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	a, _ := New(time.Minute)
	a.Add(base, 10)

	b1 := a.Buckets()
	b1[0].Lines = 999

	b2 := a.Buckets()
	if b2[0].Lines == 999 {
		t.Error("Buckets should return a copy, not a reference")
	}
}
