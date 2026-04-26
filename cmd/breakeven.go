package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/example/grpcannon/report"
)

// RunBreakeven performs a synthetic breakeven analysis using pre-collected
// step results and writes the output to stdout. In a real integration this
// would be wired to actual load-test runs at each concurrency level.
func RunBreakeven(steps []report.StepResult) error {
	if len(steps) == 0 {
		return fmt.Errorf("breakeven: no step results provided")
	}
	report.WriteBreakeven(os.Stdout, steps)
	return nil
}

// DefaultBreakevenSteps returns a set of synthetic step results useful for
// demonstration and smoke-testing the breakeven command.
func DefaultBreakevenSteps() []report.StepResult {
	return []report.StepResult{
		{Concurrency: 1, RPS: 40, P99Ms: 8, Duration: 10 * time.Second},
		{Concurrency: 5, RPS: 180, P99Ms: 15, Duration: 10 * time.Second},
		{Concurrency: 10, RPS: 320, P99Ms: 55, Duration: 10 * time.Second},
		{Concurrency: 25, RPS: 340, P99Ms: 210, Duration: 10 * time.Second},
		{Concurrency: 50, RPS: 345, P99Ms: 480, Duration: 10 * time.Second},
	}
}
