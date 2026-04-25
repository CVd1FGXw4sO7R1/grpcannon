package report

import (
	"fmt"
	"io"
	"time"
)

// ErrorRateWindow holds error rate statistics for a single time window.
type ErrorRateWindow struct {
	Start     time.Time
	End       time.Time
	Total     int
	Errors    int
	ErrorRate float64 // 0.0–1.0
}

// CalcErrorRate divides results into equal time windows and computes the
// error rate (errors / total) for each window. Windows with zero requests
// are omitted. If nWindows <= 0 it defaults to 10.
func CalcErrorRate(results []Result, nWindows int) []ErrorRateWindow {
	if len(results) == 0 {
		return nil
	}
	if nWindows <= 0 {
		nWindows = 10
	}

	var earliest, latest time.Time
	for i, r := range results {
		if i == 0 || r.Timestamp.Before(earliest) {
			earliest = r.Timestamp
		}
		if i == 0 || r.Timestamp.After(latest) {
			latest = r.Timestamp
		}
	}

	span := latest.Sub(earliest)
	if span <= 0 {
		span = time.Millisecond
	}
	windowSize := span / time.Duration(nWindows)

	type bucket struct{ total, errors int }
	buckets := make([]bucket, nWindows)

	for _, r := range results {
		idx := int(r.Timestamp.Sub(earliest) / windowSize)
		if idx >= nWindows {
			idx = nWindows - 1
		}
		buckets[idx].total++
		if r.Err != nil {
			buckets[idx].errors++
		}
	}

	var out []ErrorRateWindow
	for i, b := range buckets {
		if b.total == 0 {
			continue
		}
		start := earliest.Add(time.Duration(i) * windowSize)
		out = append(out, ErrorRateWindow{
			Start:     start,
			End:       start.Add(windowSize),
			Total:     b.total,
			Errors:    b.errors,
			ErrorRate: float64(b.errors) / float64(b.total),
		})
	}
	return out
}

// WriteErrorRate writes a human-readable error-rate-over-time table to w.
func WriteErrorRate(w io.Writer, windows []ErrorRateWindow) {
	if len(windows) == 0 {
		fmt.Fprintln(w, "no error rate data")
		return
	}
	fmt.Fprintf(w, "%-30s %8s %8s %10s\n", "window_start", "total", "errors", "error_rate")
	for _, win := range windows {
		fmt.Fprintf(w, "%-30s %8d %8d %9.2f%%\n",
			win.Start.Format(time.RFC3339),
			win.Total,
			win.Errors,
			win.ErrorRate*100,
		)
	}
}
