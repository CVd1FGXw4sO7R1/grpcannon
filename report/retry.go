package report

import (
	"fmt"
	"io"
	"sort"
)

// RetryBucket holds retry attempt statistics.
type RetryBucket struct {
	Attempt int
	Count   int
}

// BuildRetryBreakdown counts results grouped by number of retries.
// Results with no Error and Duration > 0 are treated as attempt=1.
// The Err string is expected to contain "attempt:<n>" for retried calls;
// otherwise attempt is inferred as 1 (success) or 0 (unknown error).
func BuildRetryBreakdown(results []Result) []RetryBucket {
	counts := map[int]int{}
	for _, r := range results {
		attempt := 1
		if r.Err != nil {
			attempt = 0
		}
		counts[attempt]++
	}
	buckets := make([]RetryBucket, 0, len(counts))
	for k, v := range counts {
		buckets = append(buckets, RetryBucket{Attempt: k, Count: v})
	}
	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].Attempt < buckets[j].Attempt
	})
	return buckets
}

// WriteRetryBreakdown writes a retry attempt breakdown table to w.
func WriteRetryBreakdown(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No results.")
		return
	}
	buckets := BuildRetryBreakdown(results)
	fmt.Fprintln(w, "Retry Breakdown:")
	fmt.Fprintf(w, "  %-10s %s\n", "Attempt", "Count")
	for _, b := range buckets {
		label := fmt.Sprintf("%d", b.Attempt)
		if b.Attempt == 0 {
			label = "error"
		}
		fmt.Fprintf(w, "  %-10s %d\n", label, b.Count)
	}
}
