package report

import (
	"fmt"
	"io"
	"time"
)

// WindowStat holds aggregated stats for a rolling time window.
type WindowStat struct {
	Start      time.Time
	End        time.Time
	Total      int
	Successes  int
	Failures   int
	AvgMs      float64
	P99Ms      float64
}

// CalcWindows splits results into fixed-size time windows and computes per-window stats.
func CalcWindows(results []Result, windowSize time.Duration) []WindowStat {
	if len(results) == 0 || windowSize <= 0 {
		return nil
	}

	start := results[0].StartedAt
	end := results[len(results)-1].StartedAt
	if end.Before(start) {
		return nil
	}

	var windows []WindowStat
	wStart := start
	for !wStart.After(end) {
		wEnd := wStart.Add(windowSize)
		var bucket []Result
		for _, r := range results {
			if !r.StartedAt.Before(wStart) && r.StartedAt.Before(wEnd) {
				bucket = append(bucket, r)
			}
		}
		if len(bucket) > 0 {
			stat := WindowStat{Start: wStart, End: wEnd, Total: len(bucket)}
			var durations []time.Duration
			var sumMs float64
			for _, r := range bucket {
				if r.Err == nil {
					stat.Successes++
					durations = append(durations, r.Duration)
					sumMs += float64(r.Duration.Milliseconds())
				} else {
					stat.Failures++
				}
			}
			if len(durations) > 0 {
				stat.AvgMs = sumMs / float64(len(durations))
				stat.P99Ms = Percentile(durations, 99)
			}
			windows = append(windows, stat)
		}
		wStart = wEnd
	}
	return windows
}

// WriteWindows writes per-window stats to w.
func WriteWindows(w io.Writer, windows []WindowStat) {
	if len(windows) == 0 {
		fmt.Fprintln(w, "no window data")
		return
	}
	fmt.Fprintf(w, "%-10s  %6s  %6s  %8s  %8s\n", "window", "total", "errors", "avg_ms", "p99_ms")
	for i, ws := range windows {
		fmt.Fprintf(w, "%-10d  %6d  %6d  %8.2f  %8.2f\n",
			i+1, ws.Total, ws.Failures, ws.AvgMs, ws.P99Ms)
	}
}
