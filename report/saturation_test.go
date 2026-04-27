package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeSaturationGroups() map[int][]Result {
	makeResults := func(n int, errEvery int, dur time.Duration) []Result {
		results := make([]Result, n)
		for i := range results {
			results[i] = Result{Duration: dur}
			if errEvery > 0 && i%errEvery == 0 {
				results[i].Error = errors.New("rpc error")
			}
		}
		return results
	}
	return map[int][]Result{
		1:  makeResults(100, 0, 10*time.Millisecond),
		4:  makeResults(100, 0, 20*time.Millisecond),
		8:  makeResults(100, 5, 40*time.Millisecond),
		16: makeResults(100, 2, 100*time.Millisecond),
	}
}

func TestBuildSaturation_Empty(t *testing.T) {
	sr := BuildSaturation(nil, time.Second)
	if sr == nil {
		t.Fatal("expected non-nil result")
	}
	if len(sr.Points) != 0 {
		t.Errorf("expected 0 points, got %d", len(sr.Points))
	}
}

func TestBuildSaturation_PointCount(t *testing.T) {
	groups := makeSaturationGroups()
	sr := BuildSaturation(groups, time.Second)
	if len(sr.Points) != len(groups) {
		t.Errorf("expected %d points, got %d", len(groups), len(sr.Points))
	}
}

func TestBuildSaturation_OptimalSelected(t *testing.T) {
	groups := makeSaturationGroups()
	sr := BuildSaturation(groups, time.Second)
	if sr.Optimal.Concurrency == 0 {
		t.Error("expected optimal concurrency to be set")
	}
	for _, p := range sr.Points {
		if p.Score > sr.Optimal.Score {
			t.Errorf("point with concurrency %d has higher score than optimal", p.Concurrency)
		}
	}
}

func TestBuildSaturation_ScorePositive(t *testing.T) {
	groups := map[int][]Result{
		2: {{Duration: 5 * time.Millisecond}, {Duration: 5 * time.Millisecond}},
	}
	sr := BuildSaturation(groups, time.Second)
	if sr.Optimal.Score <= 0 {
		t.Errorf("expected positive score, got %f", sr.Optimal.Score)
	}
}

func TestWriteSaturation_ValidOutput(t *testing.T) {
	groups := makeSaturationGroups()
	sr := BuildSaturation(groups, time.Second)
	var buf bytes.Buffer
	WriteSaturation(&buf, sr)
	out := buf.String()
	if !strings.Contains(out, "Saturation Analysis") {
		t.Error("expected header in output")
	}
	if !strings.Contains(out, "Optimal concurrency") {
		t.Error("expected optimal concurrency line")
	}
	if !strings.Contains(out, "*") {
		t.Error("expected optimal marker '*' in output")
	}
}

func TestWriteSaturation_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteSaturation(&buf, &SaturationResult{})
	if !strings.Contains(buf.String(), "no data") {
		t.Error("expected 'no data' message")
	}
}

func TestWriteSaturation_NilReport(t *testing.T) {
	var buf bytes.Buffer
	WriteSaturation(&buf, nil)
	if !strings.Contains(buf.String(), "no data") {
		t.Error("expected 'no data' message for nil report")
	}
}
