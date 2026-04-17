package report

import (
	"testing"
	"time"
)

func TestSummarize_Empty(t *testing.T) {
	r := Summarize(nil)
	if r.Total != 0 || r.Success != 0 {
		t.Fatal("expected zero report")
	}
}

func TestSummarize_AllSuccess(t *testing.T) {
	results := []Result{
		{Duration: ms(10)},
		{Duration: ms(20)},
		{Duration: ms(30)},
	}
	r := Summarize(results)
	if r.Total != 3 {
		t.Errorf("expected Total=3, got %d", r.Total)
	}
	if r.Success != 3 {
		t.Errorf("expected Success=3, got %d", r.Success)
	}
	if r.Failure != 0 {
		t.Errorf("expected Failure=0, got %d", r.Failure)
	}
	if r.Min != ms(10) {
		t.Errorf("expected Min=10ms, got %v", r.Min)
	}
	if r.Max != ms(30) {
		t.Errorf("expected Max=30ms, got %v", r.Max)
	}
}

func TestSummarize_WithFailures(t *testing.T) {
	results := []Result{
		{Duration: ms(10)},
		{Duration: ms(20), Err: assert_error("fail")},
	}
	r := Summarize(results)
	if r.Success != 1 || r.Failure != 1 {
		t.Errorf("expected 1 success 1 failure, got %d/%d", r.Success, r.Failure)
	}
}

func TestSummarize_Avg(t *testing.T) {
	results := []Result{
		{Duration: ms(10)},
		{Duration: ms(30)},
	}
	r := Summarize(results)
	if r.Avg != ms(20) {
		t.Errorf("expected Avg=20ms, got %v", r.Avg)
	}
}

type stubError string

func (s stubError) Error() string { return string(s) }

func assert_error(msg string) error { return stubError(msg) }

func ms(n int) time.Duration {
	return time.Duration(n) * time.Millisecond
}
