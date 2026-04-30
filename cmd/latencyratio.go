package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"grpcannon/report"
)

var (
	lrBuckets int
	lrNumerP  float64
	lrDenomP  float64
)

var latencyRatioCmd = &cobra.Command{
	Use:   "latency-ratio",
	Short: "Show P99/P50 (or custom) latency ratio across time buckets",
	Long: `Splits results into time buckets and computes the ratio of two
percentile latencies (default P99 / P50) per bucket. A rising ratio
indicates tail latency is growing faster than median latency.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunLatencyRatio(lrBuckets, lrNumerP, lrDenomP)
	},
}

func init() {
	latencyRatioCmd.Flags().IntVar(&lrBuckets, "buckets", 10, "number of time buckets")
	latencyRatioCmd.Flags().Float64Var(&lrNumerP, "numer", 99, "numerator percentile (0-100)")
	latencyRatioCmd.Flags().Float64Var(&lrDenomP, "denom", 50, "denominator percentile (0-100)")
	rootCmd.AddCommand(latencyRatioCmd)
}

// RunLatencyRatio builds synthetic results and prints the latency ratio report.
// In a real integration this would consume results from the runner.
func RunLatencyRatio(buckets int, numerP, denomP float64) error {
	if numerP <= 0 || numerP > 100 {
		return fmt.Errorf("numer must be between 1 and 100, got %.1f", numerP)
	}
	if denomP <= 0 || denomP > 100 {
		return fmt.Errorf("denom must be between 1 and 100, got %.1f", denomP)
	}

	// Placeholder: use synthetic results when no runner output is available.
	results := syntheticLatencyRatioResults()
	points := report.BuildLatencyRatio(results, buckets, numerP, denomP)
	report.WriteLatencyRatio(os.Stdout, points, numerP, denomP)
	return nil
}

func syntheticLatencyRatioResults() []report.Result {
	base := []time.Duration{
		5, 8, 10, 12, 15, 20, 25, 30, 50, 100,
	}
	results := make([]report.Result, 0, len(base)*10)
	for i := 0; i < 10; i++ {
		for _, d := range base {
			results = append(results, report.Result{
				Duration: d*time.Millisecond + time.Duration(i)*time.Millisecond,
			})
		}
	}
	return results
}
