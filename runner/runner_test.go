package runner_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/grpcannon/config"
	"github.com/grpcannon/runner"
)

func defaultConfig() *config.Config {
	cfg := config.Default()
	cfg.Target = "localhost:50051"
	cfg.Call = "pkg.Service/Method"
	cfg.Duration = 200 * time.Millisecond
	cfg.Concurrency = 2
	return cfg
}

func TestRun_Success(t *testing.T) {
	cfg := defaultConfig()
	callFn := func(ctx context.Context) error { return nil }

	r := runner.New(cfg, callFn)
	stats, err := r.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Total == 0 {
		t.Error("expected at least one call to be made")
	}
	if stats.Failures != 0 {
		t.Errorf("expected 0 failures, got %d", stats.Failures)
	}
	if stats.Successes != stats.Total {
		t.Errorf("successes (%d) != total (%d)", stats.Successes, stats.Total)
	}
}

func TestRun_WithFailures(t *testing.T) {
	cfg := defaultConfig()
	sentinelErr := errors.New("rpc error")
	callFn := func(ctx context.Context) error { return sentinelErr }

	r := runner.New(cfg, callFn)
	stats, err := r.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Failures == 0 {
		t.Error("expected failures to be recorded")
	}
	if stats.Successes != 0 {
		t.Errorf("expected 0 successes, got %d", stats.Successes)
	}
	// Failures and successes should account for all calls.
	if stats.Failures+stats.Successes != stats.Total {
		t.Errorf("failures (%d) + successes (%d) != total (%d)", stats.Failures, stats.Successes, stats.Total)
	}
}

func TestRun_InvalidConfig(t *testing.T) {
	cfg := config.Default() // missing Target and Call
	callFn := func(ctx context.Context) error { return nil }

	r := runner.New(cfg, callFn)
	_, err := r.Run(context.Background())
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestRun_DurationsRecorded(t *testing.T) {
	cfg := defaultConfig()
	callFn := func(ctx context.Context) error {
		time.Sleep(5 * time.Millisecond)
		return nil
	}

	r := runner.New(cfg, callFn)
	stats, err := r.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if int64(len(stats.Durations)) != stats.Total {
		t.Errorf("duration count %d != total %d", len(stats.Durations), stats.Total)
	}
	for _, d := range stats.Durations {
		if d < 5*time.Millisecond {
			t.Errorf("expected duration >= 5ms, got %v", d)
		}
	}
}
