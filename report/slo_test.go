package report

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func makeSLOResults(durations []time.Duration, errEvery int) []Result {
	var out []Result
	for i, d := range durations {
		var err error
		if errEvery > 0 && (i+1)%errEvery == 0 {
			err = fmt.Errorf("rpc error")
		}
		out = append(out, Result{Duration: d, Err: err})
	}
	return out
}

func TestEvaluateSLO_NilReport(t *testing.T) {
	res := EvaluateSLO(nil, SLOConfig{})
	if res.Passed {
		t.Error("expected not passed for nil report")
	}
}

func TestEvaluateSLO_AllPass(t *testing.T) {
	results := makeSLOResults([]time.Duration{
		ms(10), ms(20), ms(15), ms(12), ms(18),
	}, 0)
	r := &Report{Results: results}
	cfg := SLOConfig{
		MaxP99:         50 * time.Millisecond,
		MaxP95:         40 * time.Millisecond,
		MinSuccessRate: 1.0,
	}
	res := EvaluateSLO(r, cfg)
	if !res.Passed {
		t.Errorf("expected SLO to pass, got %+v", res)
	}
}

func TestEvaluateSLO_P99Fail(t *testing.T) {
	results := makeSLOResults([]time.Duration{
		ms(10), ms(20), ms(15), ms(200), ms(18),
	}, 0)
	r := &Report{Results: results}
	cfg := SLOConfig{MaxP99: 50 * time.Millisecond}
	res := EvaluateSLO(r, cfg)
	if res.P99Pass {
		t.Error("expected P99 to fail")
	}
	if res.Passed {
		t.Error("expected overall to fail")
	}
}

func TestEvaluateSLO_SuccessRateFail(t *testing.T) {
	results := makeSLOResults([]time.Duration{
		ms(10), ms(10), ms(10), ms(10), ms(10),
	}, 2)
	r := &Report{Results: results}
	cfg := SLOConfig{MinSuccessRate: 0.9}
	res := EvaluateSLO(r, cfg)
	if res.SuccessRatePass {
		t.Error("expected success rate to fail")
	}
}

func TestWriteSLO_Output(t *testing.T) {
	res := SLOResult{
		P99Pass: true, P95Pass: true, SuccessRatePass: false,
		P99: 30 * time.Millisecond, P95: 20 * time.Millisecond,
		SuccessRate: 0.85, Passed: false,
	}
	var buf bytes.Buffer
	WriteSLO(&buf, res)
	out := buf.String()
	if !strings.Contains(out, "FAIL") {
		t.Error("expected FAIL in output")
	}
	if !strings.Contains(out, "85.00%") {
		t.Error("expected success rate in output")
	}
}
