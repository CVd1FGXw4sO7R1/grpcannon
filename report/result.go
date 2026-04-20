package report

import "time"

// Result holds the outcome of a single gRPC call.
type Result struct {
	Timestamp time.Time
	Duration  time.Duration
	Error     error
	Attempts  int    // number of send attempts (used by retry analysis)
	StatusCode string // gRPC status code string, e.g. "OK", "UNAVAILABLE"
}

// IsSuccess returns true when the result carries no error.
func (r Result) IsSuccess() bool {
	return r.Error == nil
}
