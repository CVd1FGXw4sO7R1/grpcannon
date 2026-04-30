package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/yourusername/grpcannon/report"
)

// RunOverhead runs a synthetic overhead analysis using the provided results
// and writes a summary to stdout. It is intended to be wired into the CLI.
func RunOverhead(results []report.Result) error {
	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "overhead: no results provided")
		return nil
	}

	stats := report.CalcOverhead(results)
	report.WriteOverhead(os.Stdout, stats)
	return nil
}

// SyntheticOverheadResults builds a minimal slice of Results for demo/testing
// purposes, simulating a mix of fast and slow requests.
func SyntheticOverheadResults(n int) []report.Result {
	results := make([]report.Result, n)
	for i := 0; i < n; i++ {
		d := time.Duration((i+1)*5) * time.Millisecond
		var err error
		if i%10 == 0 {
			err = fmt.Errorf("simulated error at index %d", i)
		}
		results[i] = report.Result{Duration: d, Err: err}
	}
	return results
}
