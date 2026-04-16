package report

import (
	"fmt"
	"io"
	"time"
)

// WriteTimeline writes a time-bucketed throughput timeline to w.
// Results are grouped into buckets of bucketSize duration and RPS is shown per bucket.
func WriteTimeline(w io.Writer, results []Result, bucketSize time.Duration) error {
	if results == nil {
		_, err := fmt.Fprintln(w, "no results")
		return err
	}
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "no results")
		return err
	}
	if bucketSize <= 0 {
		bucketSize = time.Second
	}

	type bucket struct {
		total  int
		errors int
	}

	buckets := make(map[int64]*bucket)
	var minTs int64
	for i, r := range results {
		ts := r.Timestamp.UnixNano() / int64(bucketSize)
		if i == 0 {
			minTs = ts
		}
		if ts < minTs {
			minTs = ts
		}
		if buckets[ts] == nil {
			buckets[ts] = &bucket{}
		}
		buckets[ts].total++
		if r.Error != nil {
			buckets[ts].errors++
		}
	}

	// find max ts
	var maxTs int64
	for ts := range buckets {
		if ts > maxTs {
			maxTs = ts
		}
	}

	fmt.Fprintf(w, "%-12s %8s %8s %8s\n", "Offset", "Total", "Errors", "RPS")
	fmt.Fprintf(w, "%-12s %8s %8s %8s\n", "------", "-----", "------", "---")

	for ts := minTs; ts <= maxTs; ts++ {
		b := buckets[ts]
		if b == nil {
			b = &bucket{}
		}
		offset := time.Duration(ts-minTs) * bucketSize
		rps := float64(b.total) / bucketSize.Seconds()
		fmt.Fprintf(w, "%-12s %8d %8d %8.1f\n", offset.Round(time.Millisecond), b.total, b.errors, rps)
	}
	return nil
}
