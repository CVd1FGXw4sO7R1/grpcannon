package report

import (
	"fmt"
	"io"
)

// Score holds a numeric performance score and its breakdown.
type Score struct {
	Total        float64
	LatencyScore float64
	ErrorScore   float64
	ThroughputScore float64
}

// CalcScore computes a 0–100 composite performance score from a report.
// Weights: latency 40%, error rate 40%, throughput 20%.
func CalcScore(r *Report) Score {
	if r == nil || len(r.Results) == 0 {
		return Score{}
	}

	s := Summarize(r.Results)

	// Error score: 100 when 0% errors, 0 when 100% errors.
	errorScore := s.SuccessRate * 100.0

	// Latency score: based on P99. 100 if <=10ms, 0 if >=2000ms.
	p99 := Percentile(SortedDurationsMs(r.Results), 99)
	latencyScore := 100.0 - clamp((p99-10.0)/1990.0*100.0, 0, 100)

	// Throughput score: 100 if >=1000 RPS, scales down linearly.
	rpsData := CalcRPS(r.Results, 1)
	tputScore := clamp(avgRPS(rpsData)/1000.0*100.0, 0, 100)

	total := latencyScore*0.4 + errorScore*0.4 + tputScore*0.2

	return Score{
		Total:           clamp(total, 0, 100),
		LatencyScore:    latencyScore,
		ErrorScore:      errorScore,
		ThroughputScore: tputScore,
	}
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// WriteScore writes the composite score to w.
func WriteScore(w io.Writer, r *Report) {
	if r == nil || len(r.Results) == 0 {
		fmt.Fprintln(w, "Score: N/A (no results)")
		return
	}
	sc := CalcScore(r)
	fmt.Fprintf(w, "Performance Score: %.1f/100\n", sc.Total)
	fmt.Fprintf(w, "  Latency   : %.1f/100 (40%%)\n", sc.LatencyScore)
	fmt.Fprintf(w, "  Error Rate: %.1f/100 (40%%)\n", sc.ErrorScore)
	fmt.Fprintf(w, "  Throughput: %.1f/100 (20%%)\n", sc.ThroughputScore)
}
