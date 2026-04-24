package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// StepLoadBucket holds aggregated metrics for a concurrency step.
type StepLoadBucket struct {
	Concurrency int
	Total       int
	Successes   int
	Failures    int
	AvgMs       float64
	P99Ms       float64
}

// BuildStepLoad partitions results into buckets by concurrency level.
// Each result's Concurrency field is used to group entries.
// Results with zero concurrency are assigned to bucket 1.
func BuildStepLoad(results []Result, steps int) []StepLoadBucket {
	if len(results) == 0 || steps <= 0 {
		return nil
	}

	// Determine min/max concurrency from results
	minC, maxC := results[0].Concurrency, results[0].Concurrency
	for _, r := range results {
		if r.Concurrency < minC {
			minC = r.Concurrency
		}
		if r.Concurrency > maxC {
			maxC = r.Concurrency
		}
	}
	if minC <= 0 {
		minC = 1
	}
	if maxC <= 0 {
		maxC = 1
	}

	span := maxC - minC
	bucketMap := make(map[int][]time.Duration)

	for _, r := range results {
		c := r.Concurrency
		if c <= 0 {
			c = minC
		}
		var idx int
		if span == 0 {
			idx = 0
		} else {
			idx = (c - minC) * steps / (span + 1)
		}
		if idx >= steps {
			idx = steps - 1
		}
		if r.Err == nil {
			bucketMap[idx] = append(bucketMap[idx], r.Duration)
		} else {
			bucketMap[idx] = append(bucketMap[idx], 0)
		}
	}

	buckets := make([]StepLoadBucket, steps)
	for i := 0; i < steps; i++ {
		cLevel := minC + i*(span+1)/steps
		durations := bucketMap[i]
		total := len(durations)
		var successes int
		var validDurs []time.Duration
		for _, d := range durations {
			if d > 0 {
				successes++
				validDurs = append(validDurs, d)
			}
		}
		var avgMs, p99Ms float64
		if len(validDurs) > 0 {
			var sum time.Duration
			for _, d := range validDurs {
				sum += d
			}
			avgMs = float64(sum.Milliseconds()) / float64(len(validDurs))
			sorted := make([]time.Duration, len(validDurs))
			copy(sorted, validDurs)
			sort.Slice(sorted, func(a, b int) bool { return sorted[a] < sorted[b] })
			p99Ms = Percentile(sorted, 99)
		}
		buckets[i] = StepLoadBucket{
			Concurrency: cLevel,
			Total:       total,
			Successes:   successes,
			Failures:    total - successes,
			AvgMs:       avgMs,
			P99Ms:       p99Ms,
		}
	}
	return buckets
}

// WriteStepLoad writes a step-load analysis table to w.
func WriteStepLoad(w io.Writer, buckets []StepLoadBucket) {
	if len(buckets) == 0 {
		fmt.Fprintln(w, "no step-load data")
		return
	}
	fmt.Fprintf(w, "%-12s  %8s  %8s  %8s  %10s  %10s\n",
		"Concurrency", "Total", "Success", "Failure", "Avg(ms)", "P99(ms)")
	for _, b := range buckets {
		fmt.Fprintf(w, "%-12d  %8d  %8d  %8d  %10.2f  %10.2f\n",
			b.Concurrency, b.Total, b.Successes, b.Failures, b.AvgMs, b.P99Ms)
	}
}
