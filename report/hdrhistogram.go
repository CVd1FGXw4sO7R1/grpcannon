package report

import (
	"fmt"
	"io"
	"math"
	"time"
)

// HDRBucket represents a single HDR histogram bucket.
type HDRBucket struct {
	LowerMs float64
	UpperMs float64
	Count    int
	CumPct   float64
}

// BuildHDRHistogram builds an HDR-style histogram with logarithmically
// spaced buckets covering the range of observed latencies.
func BuildHDRHistogram(results []Result, numBuckets int) []HDRBucket {
	if len(results) == 0 || numBuckets <= 0 {
		return nil
	}

	var durations []time.Duration
	for _, r := range results {
		if r.Err == nil {
			durations = append(durations, r.Duration)
		}
	}
	if len(durations) == 0 {
		return nil
	}

	minMs := durations[0].Seconds() * 1000
	maxMs := durations[0].Seconds() * 1000
	for _, d := range durations[1:] {
		v := d.Seconds() * 1000
		if v < minMs {
			minMs = v
		}
		if v > maxMs {
			maxMs = v
		}
	}

	if minMs <= 0 {
		minMs = 0.001
	}
	if maxMs <= minMs {
		maxMs = minMs + 1
	}

	logMin := math.Log10(minMs)
	logMax := math.Log10(maxMs)
	step := (logMax - logMin) / float64(numBuckets)

	buckets := make([]HDRBucket, numBuckets)
	for i := 0; i < numBuckets; i++ {
		buckets[i].LowerMs = math.Pow(10, logMin+float64(i)*step)
		buckets[i].UpperMs = math.Pow(10, logMin+float64(i+1)*step)
	}

	for _, d := range durations {
		v := d.Seconds() * 1000
		for i := 0; i < numBuckets; i++ {
			if v >= buckets[i].LowerMs && (v < buckets[i].UpperMs || i == numBuckets-1) {
				buckets[i].Count++
				break
			}
		}
	}

	total := len(durations)
	cum := 0
	for i := range buckets {
		cum += buckets[i].Count
		buckets[i].CumPct = float64(cum) / float64(total) * 100
	}

	return buckets
}

// WriteHDRHistogram writes an HDR histogram to w.
func WriteHDRHistogram(w io.Writer, results []Result, numBuckets int) {
	if numBuckets <= 0 {
		numBuckets = 10
	}
	buckets := BuildHDRHistogram(results, numBuckets)
	if len(buckets) == 0 {
		fmt.Fprintln(w, "HDR Histogram: no data")
		return
	}
	fmt.Fprintln(w, "HDR Histogram (log-scale buckets):")
	fmt.Fprintf(w, "  %-12s %-12s %8s %10s\n", "Lower(ms)", "Upper(ms)", "Count", "CumPct")
	for _, b := range buckets {
		fmt.Fprintf(w, "  %-12.3f %-12.3f %8d %9.2f%%\n",
			b.LowerMs, b.UpperMs, b.Count, b.CumPct)
	}
}
