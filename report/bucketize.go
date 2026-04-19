package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Bucket holds a latency range and the count of results within it.
type Bucket struct {
	Low   time.Duration
	High  time.Duration
	Count int
}

// Bucketize partitions results into n equal-width latency buckets.
// Results with errors are excluded.
func Bucketize(results []Result, n int) []Bucket {
	if len(results) == 0 || n <= 0 {
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

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })

	min := durations[0]
	max := durations[len(durations)-1]
	if min == max {
		return []Bucket{{Low: min, High: max, Count: len(durations)}}
	}

	width := (max - min) / time.Duration(n)
	if width == 0 {
		width = 1
	}

	buckets := make([]Bucket, n)
	for i := 0; i < n; i++ {
		buckets[i] = Bucket{
			Low:  min + time.Duration(i)*width,
			High: min + time.Duration(i+1)*width,
		}
	}
	buckets[n-1].High = max + 1

	for _, d := range durations {
		for i := n - 1; i >= 0; i-- {
			if d >= buckets[i].Low {
				buckets[i].Count++
				break
			}
		}
	}

	return buckets
}

// WriteBucketize writes a bucket distribution table to w.
func WriteBucketize(w io.Writer, results []Result, n int) {
	buckets := Bucketize(results, n)
	if len(buckets) == 0 {
		fmt.Fprintln(w, "no data")
		return
	}
	fmt.Fprintf(w, "%-20s %-20s %s\n", "Low (ms)", "High (ms)", "Count")
	for _, b := range buckets {
		fmt.Fprintf(w, "%-20.2f %-20.2f %d\n",
			float64(b.Low.Microseconds())/1000.0,
			float64(b.High.Microseconds())/1000.0,
			b.Count,
		)
	}
}
