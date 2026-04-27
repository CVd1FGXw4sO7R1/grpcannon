package cmd

import (
	"fmt"
	"os"

	"github.com/nickheyer/grpcannon/report"
	"github.com/nickheyer/grpcannon/runner"
)

// DefaultLatencySegments is the default number of time segments.
const DefaultLatencySegments = 5

// RunLatencySegments runs the load test and prints a latency-segment breakdown.
func RunLatencySegments(cfg interface{ Validate() error }, segments int) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	if segments <= 0 {
		segments = DefaultLatencySegments
	}

	r, err := runner.New(cfg)
	if err != nil {
		return fmt.Errorf("runner init: %w", err)
	}

	results, err := r.Run()
	if err != nil {
		return fmt.Errorf("run failed: %w", err)
	}

	segs := report.BuildLatencySegments(results, segments)
	report.WriteLatencySegments(os.Stdout, segs)
	return nil
}
