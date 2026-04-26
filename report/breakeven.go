package report

import (
	"fmt"
	"io"
	"time"
)

// BreakevenPoint represents the concurrency level at which throughput gains
// no longer offset the increase in p99 latency.
type BreakevenPoint struct {
	Concurrency int
	RPS         float64
	P99Ms       float64
	Score       float64 // higher is better: RPS / P99Ms
}

// FindBreakeven analyses results across concurrency steps and returns the
// concurrency level that maximises RPS-to-latency efficiency.
func FindBreakeven(steps []StepResult) ([]BreakevenPoint, int) {
	if len(steps) == 0 {
		return nil, 0
	}

	points := make([]BreakevenPoint, 0, len(steps))
	best := -1
	bestScore := -1.0

	for _, s := range steps {
		if s.P99Ms <= 0 {
			continue
		}
		score := s.RPS / s.P99Ms
		points = append(points, BreakevenPoint{
			Concurrency: s.Concurrency,
			RPS:         s.RPS,
			P99Ms:       s.P99Ms,
			Score:       score,
		})
		if score > bestScore {
			bestScore = score
			best = s.Concurrency
		}
	}

	return points, best
}

// WriteBreakeven writes a human-readable breakeven analysis to w.
func WriteBreakeven(w io.Writer, steps []StepResult) {
	points, best := FindBreakeven(steps)
	if len(points) == 0 {
		fmt.Fprintln(w, "breakeven: no data")
		return
	}

	fmt.Fprintln(w, "=== Breakeven Analysis ===")
	fmt.Fprintf(w, "%-12s  %-10s  %-10s  %-10s\n", "Concurrency", "RPS", "P99 (ms)", "Score")
	for _, p := range points {
		marker := ""
		if p.Concurrency == best {
			marker = " <-- optimal"
		}
		fmt.Fprintf(w, "%-12d  %-10.2f  %-10.2f  %-10.4f%s\n",
			p.Concurrency, p.RPS, p.P99Ms, p.Score, marker)
	}
	fmt.Fprintf(w, "\nOptimal concurrency: %d\n", best)
}

// StepResult holds aggregated metrics for a single concurrency step.
type StepResult struct {
	Concurrency int
	RPS         float64
	P99Ms       float64
	Duration    time.Duration
}
