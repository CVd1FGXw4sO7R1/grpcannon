package report

import (
	"fmt"
	"io"
	"time"
)

// ConcurrencySnapshot records active workers at a point in time.
type ConcurrencySnapshot struct {
	At      time.Time
	Workers int
}

// ConcurrencySeries holds snapshots over the test duration.
type ConcurrencySeries []ConcurrencySnapshot

// Peak returns the maximum concurrency observed.
func (cs ConcurrencySeries) Peak() int {
	peak := 0
	for _, s := range cs {
		if s.Workers > peak {
			peak = s.Workers
		}
	}
	return peak
}

// Average returns the mean concurrency across all snapshots.
func (cs ConcurrencySeries) Average() float64 {
	if len(cs) == 0 {
		return 0
	}
	sum := 0
	for _, s := range cs {
		sum += s.Workers
	}
	return float64(sum) / float64(len(cs))
}

// WriteConcurrency writes a concurrency-over-time chart to w.
func WriteConcurrency(w io.Writer, series ConcurrencySeries) error {
	if len(series) == 0 {
		_, err := fmt.Fprintln(w, "No concurrency data recorded.")
		return err
	}

	peak := series.Peak()
	width := 40

	fmt.Fprintf(w, "Concurrency Over Time (peak: %d workers)\n", peak)
	fmt.Fprintf(w, "%-12s %s\n", "Elapsed", "Workers")

	start := series[0].At
	for _, s := range series {
		elapsed := s.At.Sub(start).Round(time.Millisecond)
		bar := 0
		if peak > 0 {
			bar = int(float64(s.Workers) / float64(peak) * float64(width))
		}
		fmt.Fprintf(w, "%-12s |%-*s| %d\n", elapsed, width, repeat('█', bar), s.Workers)
	}
	return nil
}

func repeat(ch rune, n int) string {
	if n <= 0 {
		return ""
	}
	runes := make([]rune, n)
	for i := range runes {
		runes[i] = ch
	}
	return string(runes)
}
