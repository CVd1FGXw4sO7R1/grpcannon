package config

import (
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Concurrency != 10 {
		t.Errorf("expected default concurrency 10, got %d", cfg.Concurrency)
	}
	if cfg.TotalRequests != 200 {
		t.Errorf("expected default total_requests 200, got %d", cfg.TotalRequests)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", cfg.Timeout)
	}
	if !cfg.Insecure {
		t.Error("expected default insecure to be true")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := Default()
	cfg.Target = "localhost:50051"
	cfg.Call = "helloworld.Greeter/SayHello"

	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_MissingTarget(t *testing.T) {
	cfg := Default()
	cfg.Call = "helloworld.Greeter/SayHello"

	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing target, got nil")
	}
}

func TestValidate_MissingCall(t *testing.T) {
	cfg := Default()
	cfg.Target = "localhost:50051"

	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing call, got nil")
	}
}

func TestValidate_ZeroConcurrency(t *testing.T) {
	cfg := Default()
	cfg.Target = "localhost:50051"
	cfg.Call = "svc/Method"
	cfg.Concurrency = 0

	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero concurrency, got nil")
	}
}

func TestValidate_ConcurrencyExceedsTotal(t *testing.T) {
	cfg := Default()
	cfg.Target = "localhost:50051"
	cfg.Call = "svc/Method"
	cfg.Concurrency = 50
	cfg.TotalRequests = 10

	if err := cfg.Validate(); err == nil {
		t.Error("expected error when concurrency exceeds total_requests, got nil")
	}
}

func TestValidate_ZeroTimeout(t *testing.T) {
	cfg := Default()
	cfg.Target = "localhost:50051"
	cfg.Call = "svc/Method"
	cfg.Timeout = 0

	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero timeout, got nil")
	}
}
