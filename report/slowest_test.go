package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeSlowestResults() []Result {
	return []Result{
		{Duration: ms(10), Error: nil},
		{Duration: ms(50), Error: nil},
		{Duration: ms(30), Error: errors.New("timeout")},
		{Duration: ms(80), Error: nil},
		{Duration: ms(5), Error: nil},
	}
}

func TestTopSlowest_Empty(t *testing.T) {
	result := TopSlowest(nil, 3)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestTopSlowest_ZeroN(t *testing.T) {
	result := TopSlowest(makeSlowestResults(), 0)
	if result != nil {
		t.Errorf("expected nil for n=0")
	}
}

func TestTopSlowest_Order(t *testing.T) {
	results := makeSlowestResults()
	top := TopSlowest(results, 3)
	if len(top) != 3 {
		t.Fatalf("expected 3, got %d", len(top))
	}
	if top[0].Duration != ms(80) {
		t.Errorf("expected 80ms first, got %v", top[0].Duration)
	}
	if top[1].Duration != ms(50) {
		t.Errorf("expected 50ms second, got %v", top[1].Duration)
	}
	if top[2].Duration != ms(30) {
		t.Errorf("expected 30ms third, got %v", top[2].Duration)
	}
}

func TestTopSlowest_NExceedsLen(t *testing.T) {
	top := TopSlowest(makeSlowestResults(), 100)
	if len(top) != 5 {
		t.Errorf("expected 5, got %d", len(top))
	}
}

func TestTopSlowest_ErrorPropagated(t *testing.T) {
	top := TopSlowest(makeSlowestResults(), 3)
	if top[2].Error != "timeout" {
		t.Errorf("expected timeout error, got %q", top[2].Error)
	}
}

func TestWriteSlowest_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	WriteSlowest(&buf, makeSlowestResults(), 3)
	out := buf.String()
	if !strings.Contains(out, "Top 3 Slowest") {
		t.Errorf("missing header: %s", out)
	}
	if !strings.Contains(out, "80ms") {
		t.Errorf("missing 80ms entry: %s", out)
	}
}

func TestWriteSlowest_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteSlowest(&buf, []Result{}, 5)
	if !strings.Contains(buf.String(), "No results") {
		t.Errorf("expected no results message, got: %s", buf.String())
	}
}

func TestTopSlowest_Ranks(t *testing.T) {
	top := TopSlowest(makeSlowestResults(), 2)
	for i, r := range top {
		if r.Rank != i+1 {
			t.Errorf("expected rank %d, got %d", i+1, r.Rank)
		}
	}
}

var _ = time.Millisecond // ensure time import used
