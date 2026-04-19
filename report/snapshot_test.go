package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeSnapshotResults() []Result {
	return []Result{
		{Duration: ms(10), Err: nil, StartedAt: time.Now()},
		{Duration: ms(20), Err: nil, StartedAt: time.Now()},
		{Duration: ms(200), Err: nil, StartedAt: time.Now()},
		{Duration: ms(5), Err: fmt.Errorf("fail"), StartedAt: time.Now()},
	}
}

func TestTakeSnapshot_Nil(t *testing.T) {
	snap := TakeSnapshot(nil)
	if snap.Total != 0 {
		t.Errorf("expected 0 total, got %d", snap.Total)
	}
}

func TestTakeSnapshot_Counts(t *testing.T) {
	r := &Report{Results: makeSnapshotResults()}
	snap := TakeSnapshot(r)
	if snap.Total != 4 {
		t.Errorf("expected total 4, got %d", snap.Total)
	}
	if snap.Successes != 3 {
		t.Errorf("expected 3 successes, got %d", snap.Successes)
	}
	if snap.Failures != 1 {
		t.Errorf("expected 1 failure, got %d", snap.Failures)
	}
}

func TestTakeSnapshot_P99Positive(t *testing.T) {
	r := &Report{Results: makeSnapshotResults()}
	snap := TakeSnapshot(r)
	if snap.P99Ms <= 0 {
		t.Errorf("expected positive P99, got %.2f", snap.P99Ms)
	}
}

func TestTakeSnapshot_Timestamp(t *testing.T) {
	before := time.Now()
	r := &Report{Results: makeSnapshotResults()}
	snap := TakeSnapshot(r)
	if snap.Timestamp.Before(before) {
		t.Error("snapshot timestamp should be after test start")
	}
}

func TestWriteSnapshot_ContainsFields(t *testing.T) {
	r := &Report{Results: makeSnapshotResults()}
	snap := TakeSnapshot(r)
	var buf bytes.Buffer
	WriteSnapshot(&buf, snap)
	out := buf.String()
	for _, want := range []string{"Snapshot", "Total", "Successes", "Failures", "Avg", "P99", "RPS"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestWriteSnapshot_EmptyReport(t *testing.T) {
	snap := TakeSnapshot(nil)
	var buf bytes.Buffer
	WriteSnapshot(&buf, snap)
	if buf.Len() == 0 {
		t.Error("expected non-empty output even for empty snapshot")
	}
}
