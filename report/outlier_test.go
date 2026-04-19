package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeOutlierResults() []Result {
	base := []time.Duration{
		10, 11, 12, 10, 13, 11, 12, 10, 11, 12,
	}
	results := make([]Result, len(base))
	for i, d := range base {
		results[i] = Result{Duration: d * time.Millisecond}
	}
	// add a clear outlier
	results = append(results, Result{Duration: 500 * time.Millisecond})
	return results
}

func TestFindOutliers_Empty(t *testing.T) {
	out := FindOutliers(nil, 2.0)
	if out != nil {
		t.Errorf("expected nil, got %v", out)
	}
}

func TestFindOutliers_Single(t *testing.T) {
	out := FindOutliers([]Result{{Duration: ms(10)}}, 2.0)
	if out != nil {
		t.Errorf("expected nil for single result")
	}
}

func TestFindOutliers_DetectsOutlier(t *testing.T) {
	results := makeOutlierResults()
	out := FindOutliers(results, 2.0)
	if len(out) == 0 {
		t.Fatal("expected at least one outlier")
	}
	if out[0].Duration != 500*time.Millisecond {
		t.Errorf("expected 500ms outlier, got %v", out[0].Duration)
	}
}

func TestFindOutliers_SkipsErrors(t *testing.T) {
	results := []Result{
		{Duration: ms(10)},
		{Duration: ms(11)},
		{Duration: ms(1000), Error: errors.New("rpc error")},
	}
	out := FindOutliers(results, 1.0)
	for _, o := range out {
		if o.Error != nil {
			t.Errorf("outlier should not include error results")
		}
	}
}

func TestFindOutliers_DefaultThreshold(t *testing.T) {
	results := makeOutlierResults()
	outDefault := FindOutliers(results, 0)
	outExplicit := FindOutliers(results, 2.0)
	if len(outDefault) != len(outExplicit) {
		t.Errorf("default threshold should equal 2.0, got %d vs %d", len(outDefault), len(outExplicit))
	}
}

func TestWriteOutliers_NoOutliers(t *testing.T) {
	var buf bytes.Buffer
	WriteOutliers(&buf, []Result{{Duration: ms(10)}, {Duration: ms(11)}}, 5.0)
	if !strings.Contains(buf.String(), "No outliers") {
		t.Errorf("expected no-outliers message, got: %s", buf.String())
	}
}

func TestWriteOutliers_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	WriteOutliers(&buf, makeOutlierResults(), 2.0)
	out := buf.String()
	if !strings.Contains(out, "Outliers") {
		t.Errorf("expected header, got: %s", out)
	}
	if !strings.Contains(out, "500") {
		t.Errorf("expected 500ms in output, got: %s", out)
	}
}
