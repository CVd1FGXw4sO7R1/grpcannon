package report

import (
	"fmt"
	"io"
	"time"
)

// RateWindow holds RPS and success metrics for a fixed time window.
type RateWindow struct {
	Start      time.Time
	End        time.Time
	Total      int
	Successes  int
	RPS        float64
	SuccessRPS float64
}

// CalcRateControl slices results into windows of windowSize duration and
// computes per-window RPS and success-RPS, enabling rate-limiting analysis.
func CalcRateControl(results []Result, windowSize time.Duration) []RateWindow {
	if len(results) == 0 || windowSize <= 0 {
		return nil
	}

	// Determine overall time range.
	minT := results[0].StartedAt
	maxT := results[0].EndedAt
	for _, r := range results {
		if r.StartedAt.Before(minT) {
			minT = r.StartedAt
		}
		if r.EndedAt.After(maxT) {
			maxT = r.EndedAt
		}
	}

	var windows []RateWindow
	for wStart := minT; wStart.Before(maxT); wStart = wStart.Add(windowSize) {
		wEnd := wStart.Add(windowSize)
		var total, successes int
		for _, r := range results {
			if r.StartedAt.Before(wEnd) && !r.StartedAt.Before(wStart) {
				total++
				if r.Err == nil {
					successes++
				}
			}
		}
		secs := windowSize.Seconds()
		windows = append(windows, RateWindow{
			Start:      wStart,
			End:        wEnd,
			Total:      total,
			Successes:  successes,
			RPS:        float64(total) / secs,
			SuccessRPS: float64(successes) / secs,
		})
	}
	return windows
}

// WriteRateControl writes a rate-control report table to w.
func WriteRateControl(w io.Writer, windows []RateWindow) {
	if len(windows) == 0 {
		fmt.Fprintln(w, "no rate-control data")
		return
	}
	fmt.Fprintf(w, "%-28s %-28s %8s %10s %12s %14s\n",
		"Window Start", "Window End", "Total", "Successes", "RPS", "Success RPS")
	for _, win := range windows {
		fmt.Fprintf(w, "%-28s %-28s %8d %10d %12.2f %14.2f\n",
			win.Start.Format(time.RFC3339),
			win.End.Format(time.RFC3339),
			win.Total,
			win.Successes,
			win.RPS,
			win.SuccessRPS,
		)
	}
}
