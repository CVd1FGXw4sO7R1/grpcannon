package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/nickbadlose/grpcannon/report"
)

// RunPercentileBox builds and prints a box-plot summary from synthetic results.
// In a real CLI invocation this would consume actual runner results.
func RunPercentileBox(results []report.Result) error {
	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "no results to summarise")
		return nil
	}
	bs := report.BuildBoxStats(results)
	return report.WriteBoxStats(os.Stdout, bs)
}

// SyntheticBoxResults returns a small set of fake results for demo / smoke testing.
func SyntheticBoxResults() []report.Result {
	durations := []time.Duration{
		5 * time.Millisecond,
		12 * time.Millisecond,
		18 * time.Millisecond,
		25 * time.Millisecond,
		30 * time.Millisecond,
		45 * time.Millisecond,
		60 * time.Millisecond,
		80 * time.Millisecond,
		95 * time.Millisecond,
		120 * time.Millisecond,
	}
	results := make([]report.Result, len(durations))
	for i, d := range durations {
		results[i] = report.Result{Duration: d}
	}
	return results
}
