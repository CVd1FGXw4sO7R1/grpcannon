package report

import "time"

// Result holds the outcome of a single gRPC call.
type Result struct {
	// Timestamp is when the call was initiated.
	Timestamp time.Time
	// Duration is the round-trip latency of the call.
	Duration time.Duration
	// Err is non-nil when the call failed.
	Err error
	// StatusCode is the gRPC status code string (e.g. "OK", "UNAVAILABLE").
	StatusCode string
	// Attempt is the 1-based retry attempt number (1 = first try).
	Attempt int
}

// IsSuccess returns true when the result carries no error.
func (r Result) IsSuccess() bool {
	return r.Err == nil
}
