package report

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// WriteTable writes a human-readable table of the report to w.
func WriteTable(r *Report, w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, "Metric\tValue")
	fmt.Fprintln(tw, strings.Repeat("-", 30))
	fmt.Fprintf(tw, "Total Requests\t%d\n", r.Total)
	fmt.Fprintf(tw, "Successful\t%d\n", r.Successful)
	fmt.Fprintf(tw, "Failed\t%d\n", r.Failed)

	if r.Total > 0 {
		successRate := float64(r.Successful) / float64(r.Total) * 100.0
		fmt.Fprintf(tw, "Success Rate\t%.2f%%\n", successRate)
	} else {
		fmt.Fprintf(tw, "Success Rate\tN/A\n")
	}

	fmt.Fprintf(tw, "Min Latency\t%s\n", roundDuration(r.Min))
	fmt.Fprintf(tw, "Mean Latency\t%s\n", roundDuration(r.Mean))
	fmt.Fprintf(tw, "Max Latency\t%s\n", roundDuration(r.Max))

	sorted := SortedDurationsMs(r.Durations)
	fmt.Fprintf(tw, "p50 Latency\t%.2fms\n", Percentile(sorted, 50))
	fmt.Fprintf(tw, "p90 Latency\t%.2fms\n", Percentile(sorted, 90))
	fmt.Fprintf(tw, "p95 Latency\t%.2fms\n", Percentile(sorted, 95))
	fmt.Fprintf(tw, "p99 Latency\t%.2fms\n", Percentile(sorted, 99))

	return tw.Flush()
}
