package report

import (
	"fmt"
	"io"
	"time"
)

// RPSPoint represents requests-per-second at a moment in time.
type RPSPoint struct {
	Second int
	RPS    float64
	Errors int
}

// CalcRPS buckets results into per-second RPS values.
func CalcRPS(results []Result) []RPSPoint {
	if len(results) == 0 {
		return nil
	}

	start := results[0].Start
	buckets := make(map[int]*RPSPoint)

	for _, r := range results {
		sec := int(r.Start.Sub(start).Seconds())
		if _, ok := buckets[sec]; !ok {
			buckets[sec] = &RPSPoint{Second: sec}
		}
		buckets[sec].RPS++
		if r.Error != nil {
			buckets[sec].Errors++
		}
	}

	max := 0
	for k := range buckets {
		if k > max {
			max = k
		}
	}

	points := make([]RPSPoint, max+1)
	for i := 0; i <= max; i++ {
		if p, ok := buckets[i]; ok {
			points[i] = *p
		} else {
			points[i] = RPSPoint{Second: i}
		}
	}
	return points
}

// WriteRPS writes a per-second RPS breakdown to w.
func WriteRPS(w io.Writer, results []Result) error {
	points := CalcRPS(results)
	if len(points) == 0 {
		_, err := fmt.Fprintln(w, "No results to display.")
		return err
	}

	fmt.Fprintf(w, "%-8s %-10s %s\n", "Second", "RPS", "Errors")
	fmt.Fprintf(w, "%-8s %-10s %s\n", "------", "---", "------")
	for _, p := range points {
		fmt.Fprintf(w, "%-8d %-10.2f %d\n", p.Second, p.RPS, p.Errors)
	}
	return nil
}

// peakRPS returns the highest RPS value across all buckets.
func peakRPS(points []RPSPoint) float64 {
	var peak float64
	for _, p := range points {
		if p.RPS > peak {
			peak = p.RPS
		}
	}
	return peak
}

// avgRPS returns the average RPS across all buckets.
func avgRPS(points []RPSPoint) float64 {
	if len(points) == 0 {
		return 0
	}
	var total float64
	for _, p := range points {
		total += p.RPS
	}
	return total / float64(len(points))
}

var _ = time.Second // ensure time import used
