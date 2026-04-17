package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// SlowestResult holds a single slow request entry.
type SlowestResult struct {
	Rank     int
	Duration time.Duration
	Error    string
}

// TopSlowest returns the n slowest results sorted descending by duration.
func TopSlowest(results []Result, n int) []SlowestResult {
	if len(results) == 0 || n <= 0 {
		return nil
	}

	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Duration > sorted[j].Duration
	})

	if n > len(sorted) {
		n = len(sorted)
	}

	out := make([]SlowestResult, n)
	for i := 0; i < n; i++ {
		errStr := ""
		if sorted[i].Error != nil {
			errStr = sorted[i].Error.Error()
		}
		out[i] = SlowestResult{
			Rank:     i + 1,
			Duration: sorted[i].Duration,
			Error:    errStr,
		}
	}
	return out
}

// WriteSlowest writes the top-n slowest requests to w.
func WriteSlowest(w io.Writer, results []Result, n int) {
	top := TopSlowest(results, n)
	if len(top) == 0 {
		fmt.Fprintln(w, "No results available.")
		return
	}
	fmt.Fprintf(w, "Top %d Slowest Requests:\n", len(top))
	fmt.Fprintf(w, "%-6s %-14s %s\n", "Rank", "Duration", "Error")
	for _, r := range top {
		errStr := r.Error
		if errStr == "" {
			errStr = "-"
		}
		fmt.Fprintf(w, "%-6d %-14s %s\n", r.Rank, roundDuration(r.Duration), errStr)
	}
}
