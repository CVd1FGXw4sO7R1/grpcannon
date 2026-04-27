package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/bojand/grpcannon/report"
	"github.com/bojand/grpcannon/runner"
)

// RunConcurrencyLatency executes load tests at each concurrency level in levels,
// aggregates results per level, and prints a concurrency-vs-latency table.
func RunConcurrencyLatency(levels []int, runFn func(concurrency int) ([]report.Result, error)) error {
	if len(levels) == 0 {
		return fmt.Errorf("at least one concurrency level is required")
	}

	groups := make(map[int][]report.Result, len(levels))

	for _, c := range levels {
		results, err := runFn(c)
		if err != nil {
			return fmt.Errorf("run at concurrency %d failed: %w", c, err)
		}
		groups[c] = results
	}

	points := report.BuildConcurrencyLatency(groups)
	report.WriteConcurrencyLatency(os.Stdout, points)
	return nil
}

// DefaultConcurrencyLevels returns a default sweep of concurrency levels.
func DefaultConcurrencyLevels() []int {
	return []int{1, 2, 4, 8, 16, 32}
}

// RunnerConcurrencyLatency is a convenience wrapper that uses runner.New to
// drive each concurrency level and collect results.
func RunnerConcurrencyLatency(cfg runner.Config, levels []int, duration time.Duration) error {
	return RunConcurrencyLatency(levels, func(concurrency int) ([]report.Result, error) {
		cfg.Concurrency = concurrency
		cfg.Duration = duration
		r, err := runner.New(cfg)
		if err != nil {
			return nil, err
		}
		return r.Run()
	})
}
