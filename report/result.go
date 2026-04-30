package report

import "time"

// Result holds the outcome of a single gRPC call.
type Result struct {
	Duration     time.Duration
	Err          error
	Attempt      int // retry attempt number (1-based)
	PayloadBytes int // size of the request payload in bytes
	Timestamp    time.Time
}

// IsSuccess returns true when the call completed without error.
func (r Result) IsSuccess() bool {
	return r.Err == nil
}
