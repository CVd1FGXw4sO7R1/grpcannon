package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Result holds the outcome of a single gRPC call.
type Result struct {
	Duration time.Duration
	Err      error
}

// Summary aggregates results from a load test run.
type Summary struct {
	Total      int
	Successes  int
	Failures   int
	TotalTime  time.Duration
	Durations  []time.Duration
}

// New builds a Summary from a slice of Results.
func New(results []Result, total time.Duration) *Summary {
	s := &Summary{
		Total:     len(results),
		TotalTime: total,
	}
	for _, r := range results {
		if r.Err != nil {
			s.Failures++
		} else {
			s.Successes++
			s.Durations = append(s.Durations, r.Duration)
		}
	}
	sort.Slice(s.Durations, func(i, j int) bool {
		return s.Durations[i] < s.Durations[j]
	})
	return s
}

// Percentile returns the p-th percentile latency (0–100).
func (s *Summary) Percentile(p float64) time.Duration {
	if len(s.Durations) == 0 {
		return 0
	}
	idx := int(float64(len(s.Durations)-1) * p / 100.0)
	return s.Durations[idx]
}

// Print writes a human-readable report to w.
func (s *Summary) Print(w io.Writer) {
	fmt.Fprintf(w, "Total requests : %d\n", s.Total)
	fmt.Fprintf(w, "Successes      : %d\n", s.Successes)
	fmt.Fprintf(w, "Failures       : %d\n", s.Failures)
	fmt.Fprintf(w, "Total time     : %s\n", s.TotalTime.Round(time.Millisecond))
	if len(s.Durations) > 0 {
		fmt.Fprintf(w, "Latency p50    : %s\n", s.Percentile(50).Round(time.Microsecond))
		fmt.Fprintf(w, "Latency p90    : %s\n", s.Percentile(90).Round(time.Microsecond))
		fmt.Fprintf(w, "Latency p99    : %s\n", s.Percentile(99).Round(time.Microsecond))
	}
}
