package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// BoxStats holds the five-number summary plus mean for a set of durations.
type BoxStats struct {
	Min    time.Duration
	Q1     time.Duration
	Median time.Duration
	Q3     time.Duration
	Max    time.Duration
	Mean   time.Duration
	IQR    time.Duration
}

// BuildBoxStats computes a box-plot five-number summary from successful results.
func BuildBoxStats(results []Result) *BoxStats {
	var durations []float64
	for _, r := range results {
		if r.Err == nil {
			durations = append(durations, float64(r.Duration))
		}
	}
	if len(durations) == 0 {
		return &BoxStats{}
	}
	sort.Float64s(durations)

	n := len(durations)
	min := time.Duration(durations[0])
	max := time.Duration(durations[n-1])

	median := percentileFromSorted(durations, 50)
	q1 := percentileFromSorted(durations, 25)
	q3 := percentileFromSorted(durations, 75)

	var sum float64
	for _, d := range durations {
		sum += d
	}
	mean := time.Duration(sum / float64(n))
	iqr := q3 - q1

	return &BoxStats{
		Min:    min,
		Q1:     q1,
		Median: median,
		Q3:     q3,
		Max:    max,
		Mean:   mean,
		IQR:    iqr,
	}
}

// percentileFromSorted returns the p-th percentile from a sorted float64 slice.
func percentileFromSorted(sorted []float64, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := p / 100.0 * float64(len(sorted)-1)
	lo := int(idx)
	hi := lo + 1
	if hi >= len(sorted) {
		return time.Duration(sorted[lo])
	}
	frac := idx - float64(lo)
	return time.Duration(sorted[lo]*(1-frac) + sorted[hi]*frac)
}

// WriteBoxStats writes a box-plot summary to w.
func WriteBoxStats(w io.Writer, bs *BoxStats) error {
	if bs == nil {
		_, err := fmt.Fprintln(w, "no box stats available")
		return err
	}
	_, err := fmt.Fprintf(w,
		"Box Plot Summary\n"+
			"  Min    : %s\n"+
			"  Q1     : %s\n"+
			"  Median : %s\n"+
			"  Mean   : %s\n"+
			"  Q3     : %s\n"+
			"  Max    : %s\n"+
			"  IQR    : %s\n",
		roundDuration(bs.Min),
		roundDuration(bs.Q1),
		roundDuration(bs.Median),
		roundDuration(bs.Mean),
		roundDuration(bs.Q3),
		roundDuration(bs.Max),
		roundDuration(bs.IQR),
	)
	return err
}
