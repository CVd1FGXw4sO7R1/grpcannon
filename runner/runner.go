package runner

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/grpcannon/config"
)

// Result holds the outcome of a single gRPC call.
type Result struct {
	Duration time.Duration
	Err      error
}

// Stats aggregates results from a load test run.
type Stats struct {
	Total     int64
	Successes int64
	Failures  int64
	Durations []time.Duration
}

// CallFunc is the function signature for executing a single gRPC call.
type CallFunc func(ctx context.Context) error

// Runner executes load tests based on the provided config.
type Runner struct {
	cfg *config.Config
	call CallFunc
}

// New creates a new Runner with the given config and call function.
func New(cfg *config.Config, call CallFunc) *Runner {
	return &Runner{cfg: cfg, call: call}
}

// Run executes the load test and returns aggregated stats.
func (r *Runner) Run(ctx context.Context) (*Stats, error) {
	if err := r.cfg.Validate(); err != nil {
		return nil, err
	}

	results := make(chan Result, r.cfg.Concurrency)
	var wg sync.WaitGroup
	var total, successes, failures int64

	workerCtx, cancel := context.WithTimeout(ctx, r.cfg.Duration)
	defer cancel()

	for i := 0; i < r.cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-workerCtx.Done():
					return
				default:
					start := time.Now()
					er			results <- Result{Duration: time.Since(start), Err: err}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var durations []time.Duration
	for res := range results {
		atomic.AddInt64(&total, 1)
		if res.Err != nil {
			atomic.AddInt64(&failures, 1)
		} else {
			atomic.AddInt64(&successes, 1)
		}
		durations = append(durations, res.Duration)
	}

	return &Stats{
		Total:     total,
		Successes: successes,
		Failures:  failures,
		Durations: durations,
	}, nil
}
