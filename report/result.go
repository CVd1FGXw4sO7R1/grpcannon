package report

import "time"

// Result holds the outcome of a single gRPC call.
type Result struct {
	Duration    time.Duration
	Error       error
	Timestamp   time.Time
	Concurrency int    // concurrency level at time of call
	Attempt     int    // retry attempt number (1-based)
	PayloadSize int    // request payload size in bytes
	StatusCode  string // gRPC status code string
}

// IsSuccess returns true when the result has no error.
func (r Result) IsSuccess() bool {
	return r.Error == nil
}
