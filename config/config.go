package config

import (
	"errors"
	"time"
)

// Config holds all configuration for a grpcannon load test run.
type Config struct {
	// Target is the gRPC server address (host:port).
	Target string `json:"target"`

	// Proto is the path to the .proto file describing the service.
	Proto string `json:"proto"`

	// Call is the fully-qualified gRPC method to invoke (e.g. "package.Service/Method").
	Call string `json:"call"`

	// Data is the JSON-encoded request payload sent to each RPC.
	Data string `json:"data"`

	// Concurrency is the number of concurrent workers sending requests.
	Concurrency int `json:"concurrency"`

	// TotalRequests is the total number of requests to send across all workers.
	TotalRequests int `json:"total_requests"`

	// Timeout is the per-RPC deadline.
	Timeout time.Duration `json:"timeout"`

	// Insecure disables TLS when true.
	Insecure bool `json:"insecure"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Concurrency:   10,
		TotalRequests: 200,
		Timeout:       5 * time.Second,
		Insecure:      true,
	}
}

// Validate checks that the Config contains the minimum required fields
// and that numeric values are within acceptable bounds.
func (c *Config) Validate() error {
	if c.Target == "" {
		return errors.New("config: target address is required")
	}
	if c.Call == "" {
		return errors.New("config: gRPC call (method) is required")
	}
	if c.Concurrency <= 0 {
		return errors.New("config: concurrency must be greater than zero")
	}
	if c.TotalRequests <= 0 {
		return errors.New("config: total_requests must be greater than zero")
	}
	if c.Concurrency > c.TotalRequests {
		return errors.New("config: concurrency cannot exceed total_requests")
	}
	if c.Timeout <= 0 {
		return errors.New("config: timeout must be a positive duration")
	}
	return nil
}
