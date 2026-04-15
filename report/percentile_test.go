package report

import (
	"testing"
	"time"
)

func TestPercentile_EmptySlice(t *testing.T) {
	if got := Percentile(nil, 50); got != 0 {
		t.Errorf("expected 0 for empty slice, got %v", got)
	}
}

func TestPercentile_SingleElement(t *testing.T) {
	data := []float64{42.0}
	for _, p := range []float64{0, 50, 99, 100} {
		if got := Percentile(data, p); got != 42.0 {
			t.Errorf("p%v: expected 42.0, got %v", p, got)
		}
	}
}

func TestPercentile_KnownValues(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	tests := []struct {
		p    float64
		want float64
	}{
		{0, 1},
		{100, 10},
		{50, 5.5},
		{90, 9.1},
	}
	for _, tc := range tests {
		got := Percentile(data, tc.p)
		diff := got - tc.want
		if diff < -0.01 || diff > 0.01 {
			t.Errorf("p%v: expected ~%v, got %v", tc.p, tc.want, got)
		}
	}
}

func TestSortedDurationsMs(t *testing.T) {
	results := []Result{
		{Duration: 30 * time.Millisecond},
		{Duration: 10 * time.Millisecond},
		{Duration: 20 * time.Millisecond},
	}
	got := SortedDurationsMs(results)
	if len(got) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(got))
	}
	expected := []float64{10, 20, 30}
	for i, v := range expected {
		if got[i] != v {
			t.Errorf("index %d: expected %v, got %v", i, v, got[i])
		}
	}
}

func TestSortedDurationsMs_Empty(t *testing.T) {
	got := SortedDurationsMs(nil)
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}
