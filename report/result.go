package report

import "time"

// Result holds the outcome of a single gRPC call.
type Result struct {
	Start    time.Time
	Duration time.Duration
	Err      error
}

// IsSuccess returns true when the call completed without error.
func (r Result) IsSuccess() bool {
	return r.Err == nil
}
