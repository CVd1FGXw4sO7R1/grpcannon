package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"grpcannon/report"
)

var percentileFenceThresholds string

var percentileFenceCmd = &cobra.Command{
	Use:   "percentile-fence",
	Short: "Evaluate latency percentile fences against thresholds",
	Long: `Evaluate whether latency percentiles breach configured thresholds.

Thresholds are specified as comma-separated PERCENTILE=MAX_MS pairs.
Example: --thresholds 50=20,90=50,99=200`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunPercentileFence(percentileFenceThresholds)
	},
}

func init() {
	percentileFenceCmd.Flags().StringVar(&percentileFenceThresholds, "thresholds", "50=50,90=100,99=200",
		"comma-separated percentile=max_ms pairs")
	rootCmd.AddCommand(percentileFenceCmd)
}

// RunPercentileFence parses thresholds and writes a fence report using synthetic results.
func RunPercentileFence(thresholdsStr string) error {
	thresholds, err := parseFenceThresholds(thresholdsStr)
	if err != nil {
		return fmt.Errorf("invalid thresholds: %w", err)
	}

	// In a real run, results would come from the runner.
	// Here we demonstrate with an empty report.
	r := &report.Report{}
	fr := report.BuildPercentileFence(r, thresholds)
	report.WritePercentileFence(os.Stdout, fr)
	return nil
}

// parseFenceThresholds parses "50=20,99=100" into a map[float64]float64.
func parseFenceThresholds(s string) (map[float64]float64, error) {
	out := make(map[float64]float64)
	if strings.TrimSpace(s) == "" {
		return out, nil
	}
	for _, pair := range strings.Split(s, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("bad pair %q, want PERCENTILE=MAX_MS", pair)
		}
		p, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return nil, fmt.Errorf("bad percentile %q: %w", parts[0], err)
		}
		m, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("bad max_ms %q: %w", parts[1], err)
		}
		out[p] = m
	}
	return out, nil
}
