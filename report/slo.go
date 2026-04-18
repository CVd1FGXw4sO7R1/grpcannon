package report

import (
	"fmt"
	"io"
	"time"
)

// SLOConfig defines thresholds for a service level objective check.
type SLOConfig struct {
	MaxP99    time.Duration
	MaxP95    time.Duration
	MinSuccessRate float64 // 0.0 - 1.0
}

// SLOResult holds the outcome of an SLO evaluation.
type SLOResult struct {
	P99Pass        bool
	P95Pass        bool
	SuccessRatePass bool
	P99            time.Duration
	P95            time.Duration
	SuccessRate    float64
	Passed         bool
}

// EvaluateSLO checks a Report against the given SLOConfig.
func EvaluateSLO(r *Report, cfg SLOConfig) SLOResult {
	if r == nil || len(r.Results) == 0 {
		return SLOResult{}
	}

	durations := SortedDurationsMs(r.Results)
	p99 := time.Duration(Percentile(durations, 99)) * time.Millisecond
	p95 := time.Duration(Percentile(durations, 95)) * time.Millisecond

	total := len(r.Results)
	successes := 0
	for _, res := range r.Results {
		if res.IsSuccess() {
			successes++
		}
	}
	successRate := float64(successes) / float64(total)

	p99Pass := cfg.MaxP99 == 0 || p99 <= cfg.MaxP99
	p95Pass := cfg.MaxP95 == 0 || p95 <= cfg.MaxP95
	srPass := cfg.MinSuccessRate == 0 || successRate >= cfg.MinSuccessRate

	return SLOResult{
		P99Pass:         p99Pass,
		P95Pass:         p95Pass,
		SuccessRatePass: srPass,
		P99:             p99,
		P95:             p95,
		SuccessRate:     successRate,
		Passed:          p99Pass && p95Pass && srPass,
	}
}

// WriteSLO writes a human-readable SLO evaluation to w.
func WriteSLO(w io.Writer, res SLOResult) {
	pass := func(b bool) string {
		if b {
			return "PASS"
		}
		return "FAIL"
	}
	fmt.Fprintln(w, "=== SLO Evaluation ===")
	fmt.Fprintf(w, "P99 Latency : %v [%s]\n", res.P99, pass(res.P99Pass))
	fmt.Fprintf(w, "P95 Latency : %v [%s]\n", res.P95, pass(res.P95Pass))
	fmt.Fprintf(w, "Success Rate: %.2f%% [%s]\n", res.SuccessRate*100, pass(res.SuccessRatePass))
	fmt.Fprintf(w, "Overall     : %s\n", pass(res.Passed))
}
