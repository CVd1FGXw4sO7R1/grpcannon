package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// TailLatency holds latency stats for a trailing window of results.
type TailLatency struct {
	WindowSize int
	Count      int
	P50        time.Duration
	P90        time.Duration
	P99        time.Duration
	Max        time.Duration
}

// CalcTailLatency computes latency percentiles over the last n results.
func CalcTailLatency(results []Result, n int) TailLatency {
	if len(results) == 0 || n <= 0 {
		return TailLatency{WindowSize: n}
	}
	window := results
	if len(results) > n {
		window = results[len(results)-n:]
	}
	var durations []float64
	for _, r := range window {
		if r.Error == nil {
			durations = append(durations, float64(r.Duration.Milliseconds()))
		}
	}
	if len(durations) == 0 {
		return TailLatency{WindowSize: n, Count: len(window)}
	}
	sort.Float64s(durations)
	toD := func(ms float64) time.Duration {
		return time.Duration(ms) * time.Millisecond
	}
	return TailLatency{
		WindowSize: n,
		Count:      len(window),
		P50:        toD(Percentile(durations, 50)),
		P90:        toD(Percentile(durations, 90)),
		P99:        toD(Percentile(durations, 99)),
		Max:        toD(durations[len(durations)-1]),
	}
}

// WriteTailLatency writes tail latency stats to w.
func WriteTailLatency(w io.Writer, t TailLatency) {
	if t.Count == 0 {
		fmt.Fprintln(w, "tail latency: no results")
		return
	}
	fmt.Fprintf(w, "tail latency (last %d):\n", t.WindowSize)
	fmt.Fprintf(w, "  p50: %s\n", t.P50)
	fmt.Fprintf(w, "  p90: %s\n", t.P90)
	fmt.Fprintf(w, "  p99: %s\n", t.P99)
	fmt.Fprintf(w, "  max: %s\n", t.Max)
}
