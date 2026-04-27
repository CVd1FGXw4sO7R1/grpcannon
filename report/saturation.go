package report

import (
	"fmt"
	"io"
	"time"
)

// SaturationPoint represents a single concurrency level measurement.
type SaturationPoint struct {
	Concurrency int
	RPS         float64
	P99Ms       float64
	ErrorRate   float64
	Score       float64
}

// SaturationResult holds the full saturation analysis.
type SaturationResult struct {
	Points  []SaturationPoint
	Optimal SaturationPoint
}

// BuildSaturation analyses results grouped by concurrency level to find the
// saturation point — the concurrency at which throughput-per-cost peaks before
// latency degrades unacceptably.
func BuildSaturation(groups map[int][]Result, window time.Duration) *SaturationResult {
	if len(groups) == 0 {
		return &SaturationResult{}
	}

	points := make([]SaturationPoint, 0, len(groups))
	for conc, results := range groups {
		var successes int
		var durations []time.Duration
		for _, r := range results {
			if r.Error == nil {
				successes++
				durations = append(durations, r.Duration)
			}
		}
		total := len(results)
		if total == 0 {
			continue
		}
		errRate := float64(total-successes) / float64(total)
		rps := 0.0
		if window > 0 {
			rps = float64(successes) / window.Seconds()
		}
		p99 := Percentile(durations, 99)
		p99ms := float64(p99.Milliseconds())
		// Score: RPS divided by (p99ms+1) penalised by error rate.
		score := (rps / (p99ms + 1)) * (1 - errRate)
		points = append(points, SaturationPoint{
			Concurrency: conc,
			RPS:         rps,
			P99Ms:       p99ms,
			ErrorRate:   errRate,
			Score:       score,
		})
	}

	var optimal SaturationPoint
	for _, p := range points {
		if p.Score > optimal.Score {
			optimal = p
		}
	}
	return &SaturationResult{Points: points, Optimal: optimal}
}

// WriteSaturation writes a human-readable saturation report to w.
func WriteSaturation(w io.Writer, sr *SaturationResult) {
	if sr == nil || len(sr.Points) == 0 {
		fmt.Fprintln(w, "saturation: no data")
		return
	}
	fmt.Fprintln(w, "Saturation Analysis")
	fmt.Fprintln(w, "-------------------")
	fmt.Fprintf(w, "%-12s %10s %10s %10s %10s\n", "Concurrency", "RPS", "P99(ms)", "ErrRate", "Score")
	for _, p := range sr.Points {
		marker := ""
		if p.Concurrency == sr.Optimal.Concurrency {
			marker = " *"
		}
		fmt.Fprintf(w, "%-12d %10.2f %10.2f %9.1f%% %10.4f%s\n",
			p.Concurrency, p.RPS, p.P99Ms, p.ErrorRate*100, p.Score, marker)
	}
	fmt.Fprintf(w, "\nOptimal concurrency: %d (score=%.4f)\n", sr.Optimal.Concurrency, sr.Optimal.Score)
}
