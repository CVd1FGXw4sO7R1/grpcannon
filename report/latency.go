package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// LatencyBand represents a latency range and its request count.
type LatencyBand struct {
	Low   time.Duration
	High  time.Duration
	Count int
}

// LatencyBands groups results into latency bands.
func LatencyBands(results []Result, bands int) []LatencyBand {
	if len(results) == 0 || bands <= 0 {
		return nil
	}

	durations := make([]time.Duration, 0, len(results))
	for _, r := range results {
		durations = append(durations, r.Duration)
	}
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })

	min := durations[0]
	max := durations[len(durations)-1]
	if min == max {
		return []LatencyBand{{Low: min, High: max, Count: len(durations)}}
	}

	step := (max - min) / time.Duration(bands)
	if step == 0 {
		step = 1
	}

	out := make([]LatencyBand, bands)
	for i := 0; i < bands; i++ {
		out[i] = LatencyBand{
			Low:  min + time.Duration(i)*step,
			High: min + time.Duration(i+1)*step,
		}
	}

	for _, d := range durations {
		idx := int((d - min) / step)
		if idx >= bands {
			idx = bands - 1
		}
		out[idx].Count++
	}

	return out
}

// WriteLatencyBands writes a simple latency band breakdown to w.
func WriteLatencyBands(w io.Writer, results []Result, bands int) error {
	if results == nil {
		_, err := fmt.Fprintln(w, "no results")
		return err
	}
	bs := LatencyBands(results, bands)
	if len(bs) == 0 {
		_, err := fmt.Fprintln(w, "no results")
		return err
	}
	fmt.Fprintf(w, "%-20s %-20s %s\n", "Low", "High", "Count")
	for _, b := range bs {
		_, err := fmt.Fprintf(w, "%-20s %-20s %d\n",
			b.Low.Round(time.Microsecond),
			b.High.Round(time.Microsecond),
			b.Count,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
