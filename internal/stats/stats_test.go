package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/stats"
)

func TestNew_InitializesStartTime(t *testing.T) {
	before := time.Now()
	c := stats.New()
	after := time.Now()

	if c.StartTime.Before(before) || c.StartTime.After(after) {
		t.Errorf("StartTime %v not in expected range [%v, %v]", c.StartTime, before, after)
	}
}

func TestRecordLine_CountsCorrectly(t *testing.T) {
	c := stats.New()

	c.RecordLine(100, true, false)
	c.RecordLine(200, false, false)
	c.RecordLine(50, true, true)

	if c.LinesRead != 3 {
		t.Errorf("LinesRead = %d, want 3", c.LinesRead)
	}
	if c.LinesMatched != 2 {
		t.Errorf("LinesMatched = %d, want 2", c.LinesMatched)
	}
	if c.LinesFiltered != 1 {
		t.Errorf("LinesFiltered = %d, want 1", c.LinesFiltered)
	}
	if c.BytesRead != 350 {
		t.Errorf("BytesRead = %d, want 350", c.BytesRead)
	}
}

func TestFinish_SetsEndTime(t *testing.T) {
	c := stats.New()
	if !c.EndTime.IsZero() {
		t.Error("EndTime should be zero before Finish")
	}
	before := time.Now()
	c.Finish()
	after := time.Now()

	if c.EndTime.Before(before) || c.EndTime.After(after) {
		t.Errorf("EndTime %v not in expected range", c.EndTime)
	}
}

func TestElapsed_BeforeFinish(t *testing.T) {
	c := stats.New()
	time.Sleep(5 * time.Millisecond)
	elapsed := c.Elapsed()
	if elapsed < 5*time.Millisecond {
		t.Errorf("Elapsed %v should be >= 5ms", elapsed)
	}
}

func TestWriteTo_ContainsFields(t *testing.T) {
	c := stats.New()
	c.RecordLine(512, true, false)
	c.RecordLine(256, true, true)
	c.Finish()

	var buf bytes.Buffer
	_, err := c.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"lines read: 2", "matched: 2", "filtered: 1", "bytes read: 768"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q: %s", want, out)
		}
	}
}
