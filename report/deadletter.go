package report

import (
	"fmt"
	"io"
	"sort"
)

// DeadLetterEntry represents a single failed request with its error message and latency.
type DeadLetterEntry struct {
	Index    int
	ErrorMsg string
	LatencyMs float64
}

// DeadLetterQueue holds the top-N failed requests for post-run inspection.
type DeadLetterQueue struct {
	Entries []DeadLetterEntry
	Total   int
}

// BuildDeadLetterQueue collects up to n failed results, sorted by latency descending.
func BuildDeadLetterQueue(results []Result, n int) *DeadLetterQueue {
	if len(results) == 0 || n <= 0 {
		return &DeadLetterQueue{}
	}

	var failed []DeadLetterEntry
	for i, r := range results {
		if r.Err != nil {
			failed = append(failed, DeadLetterEntry{
				Index:    i,
				ErrorMsg: r.Err.Error(),
				LatencyMs: float64(r.Duration.Milliseconds()),
			})
		}
	}

	total := len(failed)

	sort.Slice(failed, func(i, j int) bool {
		return failed[i].LatencyMs > failed[j].LatencyMs
	})

	if len(failed) > n {
		failed = failed[:n]
	}

	return &DeadLetterQueue{
		Entries: failed,
		Total:   total,
	}
}

// WriteDeadLetterQueue writes a dead-letter queue report to w.
func WriteDeadLetterQueue(w io.Writer, dlq *DeadLetterQueue) {
	if dlq == nil {
		fmt.Fprintln(w, "dead letter queue: nil")
		return
	}
	if len(dlq.Entries) == 0 {
		fmt.Fprintln(w, "dead letter queue: no failures")
		return
	}
	fmt.Fprintf(w, "Dead Letter Queue (showing %d of %d failures):\n", len(dlq.Entries), dlq.Total)
	fmt.Fprintf(w, "  %-6s  %-10s  %s\n", "Index", "Latency", "Error")
	for _, e := range dlq.Entries {
		fmt.Fprintf(w, "  %-6d  %-10.2fms  %s\n", e.Index, e.LatencyMs, e.ErrorMsg)
	}
}
