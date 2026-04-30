package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/example/grpcannon/report"
	"github.com/example/grpcannon/runner"
)

// DefaultFanoutLevels are the concurrency levels used when none are specified.
var DefaultFanoutLevels = []int{1, 2, 4, 8, 16, 32}

var fanoutCmd = &cobra.Command{
	Use:   "fanout",
	Short: "Run load at multiple concurrency levels and compare throughput",
	RunE: func(cmd *cobra.Command, args []string) error {
		levelsStr, _ := cmd.Flags().GetString("levels")
		levels, err := parseLevels(levelsStr)
		if err != nil {
			return fmt.Errorf("invalid levels: %w", err)
		}
		return RunFanout(levels)
	},
}

func init() {
	fanoutCmd.Flags().String("levels", "1,2,4,8,16,32", "comma-separated concurrency levels")
	rootCmd.AddCommand(fanoutCmd)
}

// RunFanout executes a fan-out load test across the given concurrency levels.
func RunFanout(levels []int) error {
	var all []report.Result
	for _, lvl := range levels {
		cfg := buildConfig(lvl)
		r, err := runner.New(cfg)
		if err != nil {
			return fmt.Errorf("runner init failed at concurrency %d: %w", lvl, err)
		}
		results, err := r.Run()
		if err != nil {
			return fmt.Errorf("run failed at concurrency %d: %w", lvl, err)
		}
		for i := range results {
			results[i].Concurrency = lvl
		}
		all = append(all, results...)
	}

	fr := report.BuildFanout(all, levels)
	report.WriteFanout(os.Stdout, fr)
	return nil
}

func parseLevels(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	levels := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		levels = append(levels, n)
	}
	return levels, nil
}
