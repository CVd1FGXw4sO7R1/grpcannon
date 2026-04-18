package report

import (
	"fmt"
	"io"
	"math"
	"time"
)

// Jitter holds jitter statistics derived from successive latency deltas.
type Jitter struct {
	Min    time.Duration
	Max    time.Duration
	Mean   time.Duration
	StdDev time.Duration
}

// CalcJitter computes inter-request latency variation from a slice of Results.
// It returns nil when fewer than two successful results are provided.
func CalcJitter(results []Result) *Jitter {
	var durations []float64
	for _, r := range results {
		if r.Error == nil {
			durations = append(durations, float64(r.Duration))
		}
	}
	if len(durations) < 2 {
		return nil
	}

	deltas := make([]float64, len(durations)-1)
	for i := 1; i < len(durations); i++ {
		d := durations[i] - durations[i-1]
		if d < 0 {
			d = -d
		}
		deltas[i-1] = d
	}

	var sum, minD, maxD float64
	minD = deltas[0]
	maxD = deltas[0]
	for _, d := range deltas {
		sum += d
		if d < minD {
			minD = d
		}
		if d > maxD {
			maxD = d
		}
	}
	mean := sum / float64(len(deltas))

	var variance float64
	for _, d := range deltas {
		diff := d - mean
		variance += diff * diff
	}
	variance /= float64(len(deltas))

	return &Jitter{
		Min:    time.Duration(minD),
		Max:    time.Duration(maxD),
		Mean:   time.Duration(mean),
		StdDev: time.Duration(math.Sqrt(variance)),
	}
}

// WriteJitter writes jitter statistics to w.
func WriteJitter(w io.Writer, results []Result) {
	j := CalcJitter(results)
	if j == nil {
		fmt.Fprintln(w, "Jitter: insufficient data")
		return
	}
	fmt.Fprintf(w, "Jitter:\n")
	fmt.Fprintf(w, "  Min:    %s\n", roundDuration(j.Min))
	fmt.Fprintf(w, "  Max:    %s\n", roundDuration(j.Max))
	fmt.Fprintf(w, "  Mean:   %s\n", roundDuration(j.Mean))
	fmt.Fprintf(w, "  StdDev: %s\n", roundDuration(j.StdDev))
}
