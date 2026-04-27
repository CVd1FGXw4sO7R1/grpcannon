package report

import (
	"bytes"
	"strings"
	"testing"
)

func makeCapacitySteps() []StepPoint {
	return []StepPoint{
		{Concurrency: 1, RPS: 50, P99Ms: 10},
		{Concurrency: 5, RPS: 200, P99Ms: 30},
		{Concurrency: 10, RPS: 350, P99Ms: 80},
		{Concurrency: 20, RPS: 400, P99Ms: 150},
		{Concurrency: 40, RPS: 390, P99Ms: 300},
	}
}

func TestBuildCapacityPlan_Empty(t *testing.T) {
	plan := BuildCapacityPlan(nil, 100)
	if plan == nil {
		t.Fatal("expected non-nil plan")
	}
	if len(plan.Points) != 0 {
		t.Errorf("expected 0 points, got %d", len(plan.Points))
	}
}

func TestBuildCapacityPlan_PointCount(t *testing.T) {
	steps := makeCapacitySteps()
	plan := BuildCapacityPlan(steps, 100)
	if len(plan.Points) != len(steps) {
		t.Errorf("expected %d points, got %d", len(steps), len(plan.Points))
	}
}

func TestBuildCapacityPlan_RecommendedMax(t *testing.T) {
	steps := makeCapacitySteps()
	// SLO = 100 ms: concurrencies 1,5,10 pass; 20 and 40 fail
	plan := BuildCapacityPlan(steps, 100)
	if plan.RecommendedMax != 10 {
		t.Errorf("expected recommended=10, got %d", plan.RecommendedMax)
	}
}

func TestBuildCapacityPlan_PeakRPS(t *testing.T) {
	steps := makeCapacitySteps()
	plan := BuildCapacityPlan(steps, 100)
	if plan.PeakRPS != 400 {
		t.Errorf("expected peakRPS=400, got %.2f", plan.PeakRPS)
	}
}

func TestBuildCapacityPlan_SaturationCapped(t *testing.T) {
	steps := []StepPoint{
		{Concurrency: 1, RPS: 10, P99Ms: 500},
	}
	plan := BuildCapacityPlan(steps, 100)
	if plan.Points[0].SaturationPct != 100 {
		t.Errorf("expected saturation capped at 100, got %.2f", plan.Points[0].SaturationPct)
	}
}

func TestBuildCapacityPlan_ZeroSLO(t *testing.T) {
	steps := makeCapacitySteps()
	plan := BuildCapacityPlan(steps, 0)
	// With zero SLO all points pass; recommended = last concurrency
	if plan.RecommendedMax != 40 {
		t.Errorf("expected recommended=40 when slo=0, got %d", plan.RecommendedMax)
	}
}

func TestWriteCapacityPlan_ValidOutput(t *testing.T) {
	steps := makeCapacitySteps()
	plan := BuildCapacityPlan(steps, 100)
	var buf bytes.Buffer
	WriteCapacityPlan(&buf, plan)
	out := buf.String()
	if !strings.Contains(out, "Capacity Plan") {
		t.Error("expected 'Capacity Plan' header")
	}
	if !strings.Contains(out, "recommended") {
		t.Error("expected 'recommended' marker in output")
	}
}

func TestWriteCapacityPlan_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteCapacityPlan(&buf, &CapacityPlan{})
	if !strings.Contains(buf.String(), "no data") {
		t.Error("expected 'no data' for empty plan")
	}
}

func TestWriteCapacityPlan_NilPlan(t *testing.T) {
	var buf bytes.Buffer
	WriteCapacityPlan(&buf, nil)
	if !strings.Contains(buf.String(), "no data") {
		t.Error("expected 'no data' for nil plan")
	}
}
