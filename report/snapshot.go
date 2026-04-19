package report

import (
	"fmt"
	"io"
	"time"
)

// Snapshot captures a point-in-time summary of a report for interval monitoring.
type Snapshot struct {
	Timestamp   time.Time
	Total       int
	Successes   int
	Failures    int
	AvgMs       float64
	P99Ms       float64
	RPS         float64
}

// TakeSnapshot creates a Snapshot from a Report at the current time.
func TakeSnapshot(r *Report) Snapshot {
	if r == nil || len(r.Results) == 0 {
		return Snapshot{Timestamp: time.Now()}
	}
	s := Summarize(r)
	durations := SortedDurationsMs(r.Results)
	p99 := Percentile(durations, 99)
	tput := CalcThroughput(r)
	return Snapshot{
		Timestamp: time.Now(),
		Total:     s.Total,
		Successes: s.Successes,
		Failures:  s.Failures,
		AvgMs:     s.AvgMs,
		P99Ms:     p99,
		RPS:       tput.RPS,
	}
}

// WriteSnapshot writes a snapshot in human-readable form to w.
func WriteSnapshot(w io.Writer, snap Snapshot) {
	fmt.Fprintf(w, "Snapshot @ %s\n", snap.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "  Total:     %d\n", snap.Total)
	fmt.Fprintf(w, "  Successes: %d\n", snap.Successes)
	fmt.Fprintf(w, "  Failures:  %d\n", snap.Failures)
	fmt.Fprintf(w, "  Avg:       %.2f ms\n", snap.AvgMs)
	fmt.Fprintf(w, "  P99:       %.2f ms\n", snap.P99Ms)
	fmt.Fprintf(w, "  RPS:       %.2f\n", snap.RPS)
}
