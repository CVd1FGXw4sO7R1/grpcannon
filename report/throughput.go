package report

import (
	"fmt"
	"io"
	"time"
)

// ThroughputStats holds requests-per-second metrics.
type ThroughputStats struct {
	TotalRequests int
	Successful    int
	Failed        int
	Duration      time.Duration
	RPS           float64
	SuccessRPS    float64
}

// CalcThroughput computes throughput stats from a slice of Results and total duration.
func CalcThroughput(results []Result, duration time.Duration) ThroughputStats {
	if len(results) == 0 || duration <= 0 {
		return ThroughputStats{}
	}
	sec := duration.Seconds()
	successful := 0
	for _, r := range results {
		if r.IsSuccess() {
			successful++
		}
	}
	total := len(results)
	return ThroughputStats{
		TotalRequests: total,
		Successful:    successful,
		Failed:        total - successful,
		Duration:      duration,
		RPS:           float64(total) / sec,
		SuccessRPS:    float64(successful) / sec,
	}
}

// WriteThroughput writes a human-readable throughput summary to w.
func WriteThroughput(w io.Writer, stats ThroughputStats) error {
	if stats.TotalRequests == 0 {
		_, err := fmt.Fprintln(w, "Throughput: no results")
		return err
	}
	_, err := fmt.Fprintf(w,
		"Throughput\n"+
			"  Duration:    %s\n"+
			"  Total:       %d\n"+
			"  Successful:  %d\n"+
			"  Failed:      %d\n"+
			"  RPS:         %.2f\n"+
			"  Success RPS: %.2f\n",
		stats.Duration.Round(time.Millisecond),
		stats.TotalRequests,
		stats.Successful,
		stats.Failed,
		stats.RPS,
		stats.SuccessRPS,
	)
	return err
}
