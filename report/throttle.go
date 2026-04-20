package report

import (
	"fmt"
	"io"
	"time"
)

// ThrottleWindow holds aggregated stats for a single time window during a throttled run.
type ThrottleWindow struct {
	WindowStart time.Time
	Duration    time.Duration
	Total       int
	Successes   int
	Failures    int
	AvgMs       float64
	RPS         float64
}

// CalcThrottle partitions results into fixed-size windows and computes per-window
// throughput and latency. windowSize must be > 0.
func CalcThrottle(results []Result, windowSize time.Duration) []ThrottleWindow {
	if len(results) == 0 || windowSize <= 0 {
		return nil
	}

	// Determine epoch from the earliest timestamp.
	epoch := results[0].Timestamp
	for _, r := range results[1:] {
		if r.Timestamp.Before(epoch) {
			epoch = r.Timestamp
		}
	}

	buckets := make(map[int64]*ThrottleWindow)
	for _, r := range results {
		idx := int64(r.Timestamp.Sub(epoch) / windowSize)
		w, ok := buckets[idx]
		if !ok {
			w = &ThrottleWindow{
				WindowStart: epoch.Add(time.Duration(idx) * windowSize),
				Duration:    windowSize,
			}
			buckets[idx] = w
		}
		w.Total++
		if r.Error == nil {
			w.Successes++
			w.AvgMs += float64(r.Duration.Milliseconds())
		} else {
			w.Failures++
		}
	}

	// Finalise averages and RPS.
	windows := make([]ThrottleWindow, 0, len(buckets))
	for _, w := range buckets {
		if w.Successes > 0 {
			w.AvgMs /= float64(w.Successes)
		}
		sec := w.Duration.Seconds()
		if sec > 0 {
			w.RPS = float64(w.Total) / sec
		}
		windows = append(windows, *w)
	}

	// Sort by window start time.
	for i := 1; i < len(windows); i++ {
		for j := i; j > 0 && windows[j].WindowStart.Before(windows[j-1].WindowStart); j-- {
			windows[j], windows[j-1] = windows[j-1], windows[j]
		}
	}
	return windows
}

// WriteThrottle writes a throttle analysis table to w.
func WriteThrottle(w io.Writer, windows []ThrottleWindow) {
	if len(windows) == 0 {
		fmt.Fprintln(w, "no throttle data available")
		return
	}
	fmt.Fprintf(w, "%-12s  %6s  %6s  %6s  %8s  %8s\n",
		"window", "total", "ok", "err", "avg_ms", "rps")
	fmt.Fprintf(w, "%s\n", "------------  ------  ------  ------  --------  --------")
	for _, win := range windows {
		fmt.Fprintf(w, "+%-11s  %6d  %6d  %6d  %8.2f  %8.2f\n",
			formatDuration(win.WindowStart.Sub(windows[0].WindowStart)),
			win.Total, win.Successes, win.Failures, win.AvgMs, win.RPS)
	}
}

func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%.3fs", d.Seconds())
}
