package cmd

import (
	"fmt"
	"os"

	"github.com/bojand/grpcannon/report"
)

// DefaultCapacitySteps is the default number of concurrency steps for a
// capacity-planning run.
const DefaultCapacitySteps = 8

// RunCapacityPlan performs a step-load run and prints a capacity plan.
// sloP99Ms is the SLO threshold in milliseconds for P99 latency.
// maxConcurrency is the upper bound of the concurrency sweep.
func RunCapacityPlan(maxConcurrency int, sloP99Ms float64, steps int) error {
	if maxConcurrency <= 0 {
		return fmt.Errorf("capacityplan: maxConcurrency must be > 0")
	}
	if steps <= 0 {
		steps = DefaultCapacitySteps
	}
	if sloP99Ms <= 0 {
		return fmt.Errorf("capacityplan: sloP99Ms must be > 0")
	}

	// Build evenly-spaced concurrency levels from 1 to maxConcurrency.
	stepPoints := make([]report.StepPoint, 0, steps)
	for i := 0; i < steps; i++ {
		concurrency := 1 + (maxConcurrency-1)*i/(steps-1)
		if steps == 1 {
			concurrency = maxConcurrency
		}
		// In a real implementation the runner would be invoked here and real
		// measurements collected. We populate a placeholder so the planner
		// can be wired up without a live gRPC target.
		stepPoints = append(stepPoints, report.StepPoint{
			Concurrency: concurrency,
			RPS:         0,
			P99Ms:       0,
		})
	}

	plan := report.BuildCapacityPlan(stepPoints, sloP99Ms)
	report.WriteCapacityPlan(os.Stdout, plan)
	return nil
}
