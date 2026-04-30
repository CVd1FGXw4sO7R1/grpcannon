package report

import (
	"fmt"
	"io"
	"sort"
)

// RequestSizeBucket holds stats for a range of request payload sizes.
type RequestSizeBucket struct {
	Label      string
	MinBytes   int
	MaxBytes   int
	Count      int
	Successes  int
	AvgLatency float64 // ms
	P99Latency float64 // ms
}

// BuildRequestSizeReport groups results by payload size buckets and computes
// per-bucket latency stats. PayloadBytes on each Result must be populated by
// the caller; results with PayloadBytes == 0 are placed in the first bucket.
func BuildRequestSizeReport(results []Result, buckets int) []RequestSizeBucket {
	if len(results) == 0 || buckets <= 0 {
		return nil
	}

	// find min/max payload size
	minB, maxB := results[0].PayloadBytes, results[0].PayloadBytes
	for _, r := range results[1:] {
		if r.PayloadBytes < minB {
			minB = r.PayloadBytes
		}
		if r.PayloadBytes > maxB {
			maxB = r.PayloadBytes
		}
	}

	span := maxB - minB
	if span == 0 {
		span = 1
	}
	step := (span + buckets - 1) / buckets

	type bucket struct {
		durations []float64
		successes int
		count     int
	}
	bs := make([]bucket, buckets)

	for _, r := range results {
		idx := (r.PayloadBytes - minB) / step
		if idx >= buckets {
			idx = buckets - 1
		}
		bs[idx].count++
		if r.Err == nil {
			bs[idx].successes++
			bs[idx].durations = append(bs[idx].durations, float64(r.Duration.Milliseconds()))
		}
	}

	out := make([]RequestSizeBucket, 0, buckets)
	for i, b := range bs {
		lo := minB + i*step
		hi := lo + step - 1
		var avg, p99 float64
		if len(b.durations) > 0 {
			sum := 0.0
			for _, d := range b.durations {
				sum += d
			}
			avg = sum / float64(len(b.durations))
			sorted := make([]float64, len(b.durations))
			copy(sorted, b.durations)
			sort.Float64s(sorted)
			p99 = sorted[int(float64(len(sorted))*0.99)]
		}
		out = append(out, RequestSizeBucket{
			Label:      fmt.Sprintf("%d-%d B", lo, hi),
			MinBytes:   lo,
			MaxBytes:   hi,
			Count:      b.count,
			Successes:  b.successes,
			AvgLatency: avg,
			P99Latency: p99,
		})
	}
	return out
}

// WriteRequestSizeReport writes a human-readable request-size breakdown.
func WriteRequestSizeReport(w io.Writer, buckets []RequestSizeBucket) {
	if len(buckets) == 0 {
		fmt.Fprintln(w, "no request-size data")
		return
	}
	fmt.Fprintf(w, "%-18s %6s %9s %10s %10s\n", "Size Range", "Count", "Successes", "Avg (ms)", "P99 (ms)")
	for _, b := range buckets {
		fmt.Fprintf(w, "%-18s %6d %9d %10.2f %10.2f\n",
			b.Label, b.Count, b.Successes, b.AvgLatency, b.P99Latency)
	}
}
