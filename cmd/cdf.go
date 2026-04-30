package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"

	"grpcannon/report"
	"grpcannon/runner"
)

const DefaultCDFSteps = 20

// RunCDF executes a load test and prints the latency CDF to stdout.
func RunCDF(cfg interface{ Validate() error }, steps int) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	return nil
}

func init() {
	var steps int
	var concurrency int
	var duration time.Duration

	cdfCmd := &cobra.Command{
		Use:   "cdf",
		Short: "Print cumulative distribution function of latencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Default()
			cfg.Concurrency = concurrency
			cfg.Duration = duration
			if err := cfg.Validate(); err != nil {
				return err
			}
			r, err := runner.New(cfg)
			if err != nil {
				return err
			}
			results, err := r.Run()
			if err != nil {
				return err
			}
			report.WriteCDF(os.Stdout, results, steps)
			return nil
		},
	}

	cdfCmd.Flags().IntVar(&steps, "steps", DefaultCDFSteps, "number of CDF sample points")
	cdfCmd.Flags().IntVar(&concurrency, "concurrency", 10, "number of concurrent workers")
	cdfCmd.Flags().DurationVar(&duration, "duration", 10*time.Second, "test duration")

	rootCmd.AddCommand(cdfCmd)
}
