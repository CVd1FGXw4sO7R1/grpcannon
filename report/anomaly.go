package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Anomaly represents a result that deviates significantly from the median.
type Anomaly struct {
	Index    int
	Duration time.Duration
	MedianMs float64
	ZScore   float64
}

// FindAnomalies detects results whose latency z-score exceeds threshold.
func FindAnomalies(results []Result, threshold float64) []Anomaly {
	if len(results) == 0 {
		return nil
	}

	var durations []float64
	for _, r := range results {
		if r.Err == nil {
			durations = append(durations, float64(r.Duration.Milliseconds()))
		}
	}
	if len(durations) == 0 {
		return nil
	}

	sorted := make([]float64, len(durations))
	copy(sorted, durations)
	sort.Float64s(sorted)

	median := sorted[len(sorted)/2]

	var sum, sumSq float64
	for _, v := range durations {
		sum += v
	}
	mean := sum / float64(len(durations))
	for _, v := range durations {
		d := v - mean
		sumSq += d * d
	}
	stddev := sqrtF(sumSq / float64(len(durations)))
	if stddev == 0 {
		return nil
	}

	var anomalies []Anomaly
	for i, r := range results {
		if r.Err != nil {
			continue
		}
		v := float64(r.Duration.Milliseconds())
		z := (v - mean) / stddev
		if z < 0 {
			z = -z
		}
		if z >= threshold {
			anomalies = append(anomalies, Anomaly{
				Index:    i,
				Duration: r.Duration,
				MedianMs: median,
				ZScore:   z,
			})
		}
	}
	return anomalies
}

// WriteAnomalies writes detected anomalies to w.
func WriteAnomalies(w io.Writer, results []Result, threshold float64) {
	anomalies := FindAnomalies(results, threshold)
	if len(anomalies) == 0 {
		fmt.Fprintln(w, "Anomalies: none detected")
		return
	}
	fmt.Fprintf(w, "Anomalies detected (z-score >= %.1f):\n", threshold)
	fmt.Fprintf(w, "  %-8s %-12s %-12s %s\n", "Index", "Duration(ms)", "Median(ms)", "Z-Score")
	for _, a := range anomalies {
		fmt.Fprintf(w, "  %-8d %-12.2f %-12.2f %.2f\n",
			a.Index, float64(a.Duration.Milliseconds()), a.MedianMs, a.ZScore)
	}
}
