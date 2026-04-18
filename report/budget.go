package report

import (
	"fmt"
	"io"
	"time"
)

// BudgetResult holds the outcome of a latency budget evaluation.
type BudgetResult struct {
	Budget   time.Duration
	P50      time.Duration
	P95      time.Duration
	P99      time.Duration
	P50Pass  bool
	P95Pass  bool
	P99Pass  bool
	Overall  bool
}

// EvaluateBudget checks whether p50, p95, and p99 latencies fit within a given budget.
func EvaluateBudget(r *Report, budget time.Duration) *BudgetResult {
	if r == nil || len(r.Results) == 0 {
		return &BudgetResult{Budget: budget}
	}

	durations := SortedDurationsMs(r.Results)
	p50 := time.Duration(Percentile(durations, 50)) * time.Millisecond
	p95 := time.Duration(Percentile(durations, 95)) * time.Millisecond
	p99 := time.Duration(Percentile(durations, 99)) * time.Millisecond

	p50Pass := p50 <= budget
	p95Pass := p95 <= budget
	p99Pass := p99 <= budget

	return &BudgetResult{
		Budget:  budget,
		P50:     p50,
		P95:     p95,
		P99:     p99,
		P50Pass: p50Pass,
		P95Pass: p95Pass,
		P99Pass: p99Pass,
		Overall: p50Pass && p95Pass && p99Pass,
	}
}

// WriteBudget writes the latency budget evaluation to w.
func WriteBudget(w io.Writer, br *BudgetResult) {
	if br == nil {
		fmt.Fprintln(w, "no budget result")
		return
	}

	pass := func(ok bool) string {
		if ok {
			return "PASS"
		}
		return "FAIL"
	}

	fmt.Fprintf(w, "Latency Budget: %v\n", br.Budget)
	fmt.Fprintf(w, "  P50: %v [%s]\n", br.P50, pass(br.P50Pass))
	fmt.Fprintf(w, "  P95: %v [%s]\n", br.P95, pass(br.P95Pass))
	fmt.Fprintf(w, "  P99: %v [%s]\n", br.P99, pass(br.P99Pass))
	overall := "PASS"
	if !br.Overall {
		overall = "FAIL"
	}
	fmt.Fprintf(w, "  Overall: %s\n", overall)
}
