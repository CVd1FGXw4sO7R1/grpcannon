package report

import (
	"fmt"
	"io"
	"time"
)

// LatencySegment represents a named time-range segment with aggregated latency stats.
type LatencySegment struct {
	Label  string
	Count  int
	Errors int
	Min    time.Duration
	Max    time.Duration
	Avg    time.Duration
	P95    time.Duration
	P99    time.Duration
}

// BuildLatencySegments divides results into n equal-width time segments and
// computes per-segment latency statistics. Results must be ordered by StartedAt.
func BuildLatencySegments(results []Result, n int) []LatencySegment {
	if len(results) == 0 || n <= 0 {
		return nil
	}

	start := results[0].StartedAt
	end := results[len(results)-1].StartedAt
	total := end.Sub(start)
	if total <= 0 {
		total = time.Millisecond
	}
	width := total / time.Duration(n)
	if width <= 0 {
		width = time.Millisecond
	}

	buckets := make([][]time.Duration, n)
	errors := make([]int, n)

	for _, r := range results {
		idx := int(r.StartedAt.Sub(start) / width)
		if idx >= n {
			idx = n - 1
		}
		if r.Error != nil {
			errors[idx]++
			continue
		}
		buckets[idx] = append(buckets[idx], r.Duration)
	}

	segments := make([]LatencySegment, n)
	for i := 0; i < n; i++ {
		seg := LatencySegment{
			Label:  fmt.Sprintf("seg%02d", i+1),
			Errors: errors[i],
			Count:  len(buckets[i]) + errors[i],
		}
		if len(buckets[i]) > 0 {
			sorted := SortedDurationsMs(buckets[i])
			var sum float64
			for _, v := range sorted {
				sum += v
			}
			seg.Min = time.Duration(sorted[0] * float64(time.Millisecond))
			seg.Max = time.Duration(sorted[len(sorted)-1] * float64(time.Millisecond))
			seg.Avg = time.Duration(sum/float64(len(sorted))) * time.Millisecond
			seg.P95 = time.Duration(Percentile(sorted, 95) * float64(time.Millisecond))
			seg.P99 = time.Duration(Percentile(sorted, 99) * float64(time.Millisecond))
		}
		segments[i] = seg
	}
	return segments
}

// WriteLatencySegments writes a human-readable latency segment table to w.
func WriteLatencySegments(w io.Writer, segments []LatencySegment) {
	if len(segments) == 0 {
		fmt.Fprintln(w, "no latency segment data")
		return
	}
	fmt.Fprintf(w, "%-8s %6s %6s %8s %8s %8s %8s %8s\n",
		"Segment", "Count", "Errors", "Min", "Max", "Avg", "P95", "P99")
	for _, s := range segments {
		fmt.Fprintf(w, "%-8s %6d %6d %8s %8s %8s %8s %8s\n",
			s.Label, s.Count, s.Errors,
			roundDuration(s.Min), roundDuration(s.Max),
			roundDuration(s.Avg), roundDuration(s.P95), roundDuration(s.P99))
	}
}
