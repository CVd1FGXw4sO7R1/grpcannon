package report

import (
	"fmt"
	"io"
	"sort"
)

// ErrorBreakdown holds counts of each unique error message.
type ErrorBreakdown struct {
	Errors map[string]int
	Total  int
}

// BuildErrorBreakdown counts occurrences of each error in results.
func BuildErrorBreakdown(results []Result) *ErrorBreakdown {
	counts := make(map[string]int)
	for _, r := range results {
		if r.Error != nil {
			counts[r.Error.Error()]++
		}
	}
	total := 0
	for _, v := range counts {
		total += v
	}
	return &ErrorBreakdown{Errors: counts, Total: total}
}

// WriteErrorBreakdown writes a human-readable error breakdown to w.
func WriteErrorBreakdown(w io.Writer, results []Result) error {
	bd := BuildErrorBreakdown(results)
	if bd.Total == 0 {
		_, err := fmt.Fprintln(w, "No errors recorded.")
		return err
	}

	// Sort keys for deterministic output.
	keys := make([]string, 0, len(bd.Errors))
	for k := range bd.Errors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "Error Breakdown (%d total errors):\n", bd.Total)
	for _, k := range keys {
		pct := float64(bd.Errors[k]) / float64(bd.Total) * 100
		fmt.Fprintf(w, "  %-60s %5d  (%5.1f%%)\n", k, bd.Errors[k], pct)
	}
	return nil
}
