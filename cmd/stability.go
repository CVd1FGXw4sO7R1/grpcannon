package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/nickspring/grpcannon/report"
)

// DefaultStabilityWindows is the default number of time windows for analysis.
const DefaultStabilityWindows = 10

// DefaultP99JitterThresholdMs is the default p99 jitter threshold in ms.
const DefaultP99JitterThresholdMs = 50.0

// DefaultErrorDeltaThreshold is the default error-rate delta threshold.
const DefaultErrorDeltaThreshold = 0.05

// RunStability evaluates load-test stability from a set of results and prints
// a summary to stdout. It returns an error when the run is deemed unstable.
func RunStability(
	results []report.Result,
	windows int,
	p99JitterMs float64,
	errorDelta float64,
) error {
	if windows <= 0 {
		windows = DefaultStabilityWindows
	}
	if p99JitterMs <= 0 {
		p99JitterMs = DefaultP99JitterThresholdMs
	}
	if errorDelta <= 0 {
		errorDelta = DefaultErrorDeltaThreshold
	}

	start := time.Now()
	sr := report.BuildStability(results, windows, p99JitterMs, errorDelta)
	elapsed := time.Since(start)

	report.WriteStability(os.Stdout, sr)
	fmt.Fprintf(os.Stdout, "  Analysis time  : %s\n", elapsed.Round(time.Microsecond))

	if !sr.Stable {
		return fmt.Errorf("stability check failed: %s", sr.Reason)
	}
	return nil
}
