package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// WriteFlamegraph writes a simplified text-based flamegraph-style breakdown
// of latency buckets to w.
func WriteFlamegraph(w io.Writer, results []Result) error {
	if results == nil {
		return fmt.Errorf("results is nil")
	}
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "no results")
		return err
	}

	buckets := map[string]int{
		"0-5ms":    0,
		"5-10ms":   0,
		"10-25ms":  0,
		"25-50ms":  0,
		"50-100ms": 0,
		"100ms+":   0,
	}
	order := []string{"0-5ms", "5-10ms", "10-25ms", "25-50ms", "50-100ms", "100ms+"}
	limits := []time.Duration{
		5 * time.Millisecond,
		10 * time.Millisecond,
		25 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
	}

	for _, r := range results {
		assigned := false
		for i, limit := range limits {
			if r.Duration < limit {
				buckets[order[i]]++
				assigned = true
				break
			}
		}
		if !assigned {
			buckets["100ms+"]++
		}
	}

	total := len(results)
	sort.Strings(order) // stable label order for output
	order = []string{"0-5ms", "5-10ms", "10-25ms", "25-50ms", "50-100ms", "100ms+"}

	fmt.Fprintln(w, "Flamegraph (latency distribution):")
	for _, label := range order {
		count := buckets[label]
		pct := float64(count) / float64(total) * 100
		bar := buildBar(count, total, 40)
		fmt.Fprintf(w, "  %-10s |%s| %d (%.1f%%)\n", label, bar, count, pct)
	}
	return nil
}

func buildBar(count, total, width int) string {
	if total == 0 {
		return ""
	}
	filled := int(float64(count) / float64(total) * float64(width))
	bar := make([]byte, width)
	for i := range bar {
		if i < filled {
			bar[i] = '#'
		} else {
			bar[i] = ' '
		}
	}
	return string(bar)
}
