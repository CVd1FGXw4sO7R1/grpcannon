package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPercentileRank_Empty(t *testing.T) {
	rank := PercentileRank(nil, 100*time.Millisecond)
	if rank != 0 {
		t.Errorf("expected 0, got %f", rank)
	}
}

func TestPercentileRank_AllBelow(t *testing.T) {
	results := []Result{
		{Duration: ms(10)},
		{Duration: ms(20)},
		{Duration: ms(30)},
	}
	rank := PercentileRank(results, ms(50))
	if rank != 100.0 {
		t.Errorf("expected 100.0, got %f", rank)
	}
}

func TestPercentileRank_AllAbove(t *testing.T) {
	results := []Result{
		{Duration: ms(200)},
		{Duration: ms(300)},
	}
	rank := PercentileRank(results, ms(50))
	if rank != 0.0 {
		t.Errorf("expected 0.0, got %f", rank)
	}
}

func TestPercentileRank_Partial(t *testing.T) {
	results := []Result{
		{Duration: ms(10)},
		{Duration: ms(50)},
		{Duration: ms(100)},
		{Duration: ms(200)},
	}
	rank := PercentileRank(results, ms(100))
	if rank != 75.0 {
		t.Errorf("expected 75.0, got %f", rank)
	}
}

func TestWritePercentileRank_Empty(t *testing.T) {
	var buf bytes.Buffer
	WritePercentileRank(&buf, nil)
	if !strings.Contains(buf.String(), "No results") {
		t.Errorf("expected no-results message, got: %s", buf.String())
	}
}

func TestWritePercentileRank_ValidOutput(t *testing.T) {
	results := make([]Result, 10)
	for i := range results {
		results[i] = Result{Duration: ms((i + 1) * 50)}
	}
	var buf bytes.Buffer
	WritePercentileRank(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "Percentile Ranks") {
		t.Errorf("expected header in output, got: %s", out)
	}
	if !strings.Contains(out, "500ms") {
		t.Errorf("expected 500ms threshold in output, got: %s", out)
	}
}
