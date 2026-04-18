package report

import (
	"fmt"
	"io"
)

// WarmupSummary holds stats comparing warm-up vs steady-state results.
type WarmupSummary struct {
	WarmupCount    int
	SteadyCount    int
	WarmupAvgMs    float64
	SteadyAvgMs    float64
	ImprovementPct float64
}

// CalcWarmup splits results at warmupN and computes avg latency for each phase.
func CalcWarmup(results []Result, warmupN int) WarmupSummary {
	if len(results) == 0 || warmupN <= 0 {
		return WarmupSummary{}
	}
	if warmupN > len(results) {
		warmupN = len(results)
	}
	avgMs := func(rs []Result) float64 {
		var sum float64
		var count int
		for _, r := range rs {
			if r.Error == nil {
				sum += float64(r.Duration.Milliseconds())
				count++
			}
		}
		if count == 0 {
			return 0
		}
		return sum / float64(count)
	}
	warm := results[:warmupN]
	steady := results[warmupN:]
	wAvg := avgMs(warm)
	sAvg := avgMs(steady)
	var imp float64
	if wAvg > 0 {
		imp = (wAvg - sAvg) / wAvg * 100
	}
	return WarmupSummary{
		WarmupCount:    len(warm),
		SteadyCount:    len(steady),
		WarmupAvgMs:    wAvg,
		SteadyAvgMs:    sAvg,
		ImprovementPct: imp,
	}
}

// WriteWarmup writes a warmup comparison summary to w.
func WriteWarmup(w io.Writer, s WarmupSummary) {
	if s.WarmupCount == 0 {
		fmt.Fprintln(w, "warmup: no data")
		return
	}
	fmt.Fprintf(w, "warmup vs steady-state:\n")
	fmt.Fprintf(w, "  warmup  requests: %d  avg: %.2fms\n", s.WarmupCount, s.WarmupAvgMs)
	fmt.Fprintf(w, "  steady  requests: %d  avg: %.2fms\n", s.SteadyCount, s.SteadyAvgMs)
	fmt.Fprintf(w, "  improvement: %.1f%%\n", s.ImprovementPct)
}
