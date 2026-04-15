package report

import "time"

// Summary holds pre-computed aggregate statistics for a completed run.
type Summary struct {
	Total      int
	Successes  int
	Failures   int
	SuccessRate float64
	AvgMs      float64
	MinMs      float64
	MaxMs      float64
	P50Ms      float64
	P90Ms      float64
	P99Ms      float64
	TotalTime  time.Duration
}

// Summarize computes a Summary from a slice of Results and the total wall-clock
// duration of the run.
func Summarize(results []Result, totalTime time.Duration) Summary {
	s := Summary{
		Total:     len(results),
		TotalTime: totalTime,
	}
	if len(results) == 0 {
		return s
	}

	var sumMs float64
	for _, r := range results {
		if r.Err == nil {
			s.Successes++
		} else {
			s.Failures++
		}
		sumMs += msFloat(r.Duration)
	}
	s.SuccessRate = float64(s.Successes) / float64(s.Total) * 100
	s.AvgMs = sumMs / float64(s.Total)

	sorted := SortedDurationsMs(results)
	s.MinMs = sorted[0]
	s.MaxMs = sorted[len(sorted)-1]
	s.P50Ms = Percentile(sorted, 50)
	s.P90Ms = Percentile(sorted, 90)
	s.P99Ms = Percentile(sorted, 99)
	return s
}
