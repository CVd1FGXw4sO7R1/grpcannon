package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/nickheyer/grpcannon/report"
	"github.com/nickheyer/grpcannon/runner"
	"github.com/spf13/cobra"
)

var (
	pacingTargetRPS float64
	pacingWindows   int
	pacingDuration  time.Duration
	pacingConc      int
)

var pacingCmd = &cobra.Command{
	Use:   "pacing",
	Short: "Analyse how closely actual request rate matches a target RPS",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunPacing(pacingTargetRPS, pacingWindows, pacingDuration, pacingConc)
	},
}

func init() {
	pacingCmd.Flags().Float64Var(&pacingTargetRPS, "target-rps", 10, "target requests per second")
	pacingCmd.Flags().IntVar(&pacingWindows, "windows", 5, "number of analysis windows")
	pacingCmd.Flags().DurationVar(&pacingDuration, "duration", 10*time.Second, "total load duration")
	pacingCmd.Flags().IntVar(&pacingConc, "concurrency", 1, "number of concurrent workers")
	rootCmd.AddCommand(pacingCmd)
}

// RunPacing executes a pacing analysis run and prints the result.
func RunPacing(targetRPS float64, windows int, duration time.Duration, concurrency int) error {
	cfg := defaultRunConfig()
	cfg.Duration = duration
	cfg.Concurrency = concurrency

	r, err := runner.New(cfg)
	if err != nil {
		return fmt.Errorf("runner init: %w", err)
	}

	results, err := r.Run()
	if err != nil {
		return fmt.Errorf("run failed: %w", err)
	}

	pr := report.CalcPacing(results, targetRPS, windows)
	report.WritePacing(os.Stdout, pr)
	return nil
}
