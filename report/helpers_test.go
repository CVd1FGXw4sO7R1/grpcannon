package report

import "time"

// ms returns a Duration of n milliseconds, used across report tests.
func ms(n int) time.Duration {
	return time.Duration(n) * time.Millisecond
}
