package report

import (
	"fmt"
	"io"
	"time"
)

// WriteMarkdown writes a Markdown-formatted report to the given writer.
func WriteMarkdown(w io.Writer, r *Report) error {
	if r == nil {
		return fmt.Errorf("report is nil")
	}

	fmt.Fprintln(w, "# gRPCannon Load Test Report")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "## Summary")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "| Metric | Value |\n")
	fmt.Fprintf(w, "|--------|-------|\n")
	fmt.Fprintf(w, "| Total Requests | %d |\n", r.Total)
	fmt.Fprintf(w, "| Successful | %d |\n", r.Success)
	fmt.Fprintf(w, "| Failed | %d |\n", r.Failure)

	if r.Total > 0 {
		successRate := float64(r.Success) / float64(r.Total) * 100
		fmt.Fprintf(w, "| Success Rate | %.2f%% |\n", successRate)
	} else {
		fmt.Fprintf(w, "| Success Rate | N/A |\n")
	}

	fmt.Fprintf(w, "| Fastest | %s |\n", roundDuration(r.Fastest))
	fmt.Fprintf(w, "| Slowest | %s |\n", roundDuration(r.Slowest))
	fmt.Fprintf(w, "| Average | %s |\n", roundDuration(r.Average))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Latency Percentiles")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "| Percentile | Latency |\n")
	fmt.Fprintf(w, "|------------|---------|\n")

	percentiles := []struct {
		label string
		val   time.Duration
	}{
		{"p50", r.P50},
		{"p90", r.P90},
		{"p95", r.P95},
		{"p99", r.P99},
	}

	for _, p := range percentiles {
		fmt.Fprintf(w, "| %s | %s |\n", p.label, roundDuration(p.val))
	}

	return nil
}
