package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// HistoryEntry represents a single run's percentile snapshot.
type HistoryEntry struct {
	Timestamp time.Time
	Label     string
	P50Ms     float64
	P90Ms     float64
	P99Ms     float64
	SuccessRate float64
}

// PercentileHistory holds multiple run entries for trend comparison.
type PercentileHistory struct {
	Entries []HistoryEntry
}

// BuildPercentileHistory constructs a PercentileHistory from a slice of
// labelled Reports, ordered by timestamp ascending.
func BuildPercentileHistory(runs []struct {
	Label  string
	Report *Report
}) *PercentileHistory {
	if len(runs) == 0 {
		return &PercentileHistory{}
	}

	entries := make([]HistoryEntry, 0, len(runs))
	for _, r := range runs {
		if r.Report == nil {
			continue
		}
		sorted := SortedDurationsMs(r.Report.Results)
		var p50, p90, p99 float64
		if len(sorted) > 0 {
			p50 = Percentile(sorted, 50)
			p90 = Percentile(sorted, 90)
			p99 = Percentile(sorted, 99)
		}
		total := len(r.Report.Results)
		successes := 0
		for _, res := range r.Report.Results {
			if res.Error == nil {
				successes++
			}
		}
		var sr float64
		if total > 0 {
			sr = float64(successes) / float64(total) * 100
		}
		entries = append(entries, HistoryEntry{
			Timestamp:   r.Report.Start,
			Label:       r.Label,
			P50Ms:       p50,
			P90Ms:       p90,
			P99Ms:       p99,
			SuccessRate: sr,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	return &PercentileHistory{Entries: entries}
}

// WritePercentileHistory writes a textual history table to w.
func WritePercentileHistory(w io.Writer, h *PercentileHistory) {
	if h == nil || len(h.Entries) == 0 {
		fmt.Fprintln(w, "no history entries")
		return
	}
	fmt.Fprintf(w, "%-24s %-16s %8s %8s %8s %10s\n",
		"Timestamp", "Label", "P50(ms)", "P90(ms)", "P99(ms)", "Success%")
	fmt.Fprintf(w, "%s\n", fmt.Sprintf("%0*d", 72, 0))
	for _, e := range h.Entries {
		fmt.Fprintf(w, "%-24s %-16s %8.2f %8.2f %8.2f %9.1f%%\n",
			e.Timestamp.Format(time.RFC3339),
			e.Label,
			e.P50Ms,
			e.P90Ms,
			e.P99Ms,
			e.SuccessRate,
		)
	}
}
