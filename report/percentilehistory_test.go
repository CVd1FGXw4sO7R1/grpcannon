package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeHistoryRuns(labels []string, durations [][]time.Duration, errFlags [][]bool) []struct {
	Label  string
	Report *Report
} {
	base := time.Now()
	runs := make([]struct {
		Label  string
		Report *Report
	}, len(labels))
	for i, label := range labels {
		results := make([]Result, len(durations[i]))
		for j, d := range durations[i] {
			var err error
			if errFlags[i][j] {
				err = errSentinel
			}
			results[j] = Result{Duration: d, Error: err}
		}
		runs[i] = struct {
			Label  string
			Report *Report
		}{
			Label:  label,
			Report: &Report{Results: results, Start: base.Add(time.Duration(i) * time.Minute)},
		}
	}
	return runs
}

var errSentinel = fmt.Errorf("err")

func TestBuildPercentileHistory_Empty(t *testing.T) {
	h := BuildPercentileHistory(nil)
	if len(h.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(h.Entries))
	}
}

func TestBuildPercentileHistory_EntryCount(t *testing.T) {
	runs := makeHistoryRuns(
		[]string{"run1", "run2"},
		[][]time.Duration{{ms(10), ms(20)}, {ms(30), ms(40)}},
		[][]bool{{false, false}, {false, false}},
	)
	h := BuildPercentileHistory(runs)
	if len(h.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(h.Entries))
	}
}

func TestBuildPercentileHistory_P99Positive(t *testing.T) {
	runs := makeHistoryRuns(
		[]string{"run1"},
		[][]time.Duration{{ms(10), ms(20), ms(100)}},
		[][]bool{{false, false, false}},
	)
	h := BuildPercentileHistory(runs)
	if h.Entries[0].P99Ms <= 0 {
		t.Fatalf("expected positive P99, got %f", h.Entries[0].P99Ms)
	}
}

func TestBuildPercentileHistory_SuccessRate(t *testing.T) {
	runs := makeHistoryRuns(
		[]string{"run1"},
		[][]time.Duration{{ms(10), ms(20)}},
		[][]bool{{false, true}},
	)
	h := BuildPercentileHistory(runs)
	if h.Entries[0].SuccessRate != 50.0 {
		t.Fatalf("expected 50%% success rate, got %f", h.Entries[0].SuccessRate)
	}
}

func TestBuildPercentileHistory_Sorted(t *testing.T) {
	base := time.Now()
	runs := []struct {
		Label  string
		Report *Report
	}{
		{Label: "later", Report: &Report{Results: []Result{{Duration: ms(5)}}, Start: base.Add(2 * time.Minute)}},
		{Label: "earlier", Report: &Report{Results: []Result{{Duration: ms(5)}}, Start: base}},
	}
	h := BuildPercentileHistory(runs)
	if h.Entries[0].Label != "earlier" {
		t.Fatalf("expected entries sorted by time, got %s first", h.Entries[0].Label)
	}
}

func TestWritePercentileHistory_ValidOutput(t *testing.T) {
	runs := makeHistoryRuns(
		[]string{"run1"},
		[][]time.Duration{{ms(10), ms(50), ms(200)}},
		[][]bool{{false, false, false}},
	)
	h := BuildPercentileHistory(runs)
	var buf bytes.Buffer
	WritePercentileHistory(&buf, h)
	out := buf.String()
	if !strings.Contains(out, "run1") {
		t.Fatalf("expected label in output, got: %s", out)
	}
	if !strings.Contains(out, "P99") {
		t.Fatalf("expected header in output, got: %s", out)
	}
}

func TestWritePercentileHistory_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WritePercentileHistory(&buf, &PercentileHistory{})
	if !strings.Contains(buf.String(), "no history") {
		t.Fatalf("expected empty message, got: %s", buf.String())
	}
}
