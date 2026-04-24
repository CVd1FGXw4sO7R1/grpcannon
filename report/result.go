package report

import (
	"errors"
	"time"
)

// errSentinel is a package-level sentinel error used in tests.
var errSentinel = errors.New("sentinel error")

// Result represents the outcome of a single gRPC call.
type Result struct {
	// Duration is the round-trip latency of the call.
	Duration time.Duration
	// Err is non-nil when the call failed.
	Err error
	// Timestamp is when the call was initiated.
	Timestamp time.Time
	// Concurrency is the number of concurrent workers active when this
	// call was issued. Zero means unset.
	Concurrency int
	// Attempt is the 1-based retry attempt number (1 = first try).
	Attempt int
	// StatusCode is the gRPC status code string, e.g. "OK", "UNAVAILABLE".
	StatusCode string
}

// IsSuccess reports whether the result represents a successful call.
func (r Result) IsSuccess() bool {
	return r.Err == nil
}
