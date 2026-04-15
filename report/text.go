package report

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// WriteText writes a human-readable summary report to the given writer.
func WriteText(r *Report, w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "Summary\n")
	fmt.Fprintf(tw, "-------\n")
	fmt.Fprintf(tw, "Total requests:\t%d\n", r.Total)
	fmt.Fprintf(tw, "Successful:\t%d\n", r.Success)
	fmt.Fprintf(tw, "Failed:\t%d\n", r.Failure)

	if r.Total > 0 {
		successRate := float64(r.Success) / float64(r.Total) * 100
		fmt.Fprintf(tw, "Success rate:\t%.2f%%\n", successRate)
	}

	fmt.Fprintf(tw, "\nLatency\n")
	fmt.Fprintf(tw, "-------\n")
	fmt.Fprintf(tw, "Min:\t%s\n", roundDuration(r.Fastest))
	fmt.Fprintf(tw, "Mean:\t%s\n", roundDuration(r.Average))
	fmt.Fprintf(tw, "Max:\t%s\n", roundDuration(r.Slowest))

	fmt.Fprintf(tw, "\nPercentiles\n")
	fmt.Fprintf(tw, "-----------\n")
	fmt.Fprintf(tw, "p50:\t%s\n", roundDuration(r.P50))
	fmt.Fprintf(tw, "p90:\t%s\n", roundDuration(r.P90))
	fmt.Fprintf(tw, "p95:\t%s\n", roundDuration(r.P95))
	fmt.Fprintf(tw, "p99:\t%s\n", roundDuration(r.P99))

	return tw.Flush()
}

func roundDuration(d time.Duration) string {
	return d.Round(time.Microsecond).String()
}
