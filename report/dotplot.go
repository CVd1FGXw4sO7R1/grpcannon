package report

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// WriteDotPlot writes a simple ASCII dot plot of latency distribution to w.
func WriteDotPlot(w io.Writer, r *Report) error {
	if r == nil || len(r.Results) == 0 {
		_, err := fmt.Fprintln(w, "No results to plot.")
		return err
	}

	const buckets = 10
	const width = 40

	durations := SortedDurationsMs(r.Results)
	if len(durations) == 0 {
		_, err := fmt.Fprintln(w, "No durations to plot.")
		return err
	}

	min := durations[0]
	max := durations[len(durations)-1]
	span := max - min
	if span == 0 {
		span = 1
	}

	counts := make([]int, buckets)
	for _, d := range durations {
		idx := int((d - min) / span * float64(buckets-1))
		if idx >= buckets {
			idx = buckets - 1
		}
		counts[idx]++
	}

	maxCount := 0
	for _, c := range counts {
		if c > maxCount {
			maxCount = c
		}
	}

	bucketSize := span / float64(buckets)
	fmt.Fprintln(w, "Latency Distribution (ms):")
	for i := 0; i < buckets; i++ {
		lo := min + float64(i)*bucketSize
		hi := lo + bucketSize
		bar := 0
		if maxCount > 0 {
			bar = counts[i] * width / maxCount
		}
		fmt.Fprintf(w, "  %6.2f - %6.2f ms | %s (%d)\n", lo, hi, strings.Repeat("█", bar), counts[i])
	}

	fmt.Fprintf(w, "  Min: %.2f ms  Max: %.2f ms  Mean: %.2f ms\n",
		min, max, msFloat(time.Duration(r.Mean)))
	return nil
}
