package report

import "sort"

// Percentile calculates the nth percentile (0-100) from a sorted slice of
// float64 durations (in milliseconds). Returns 0 if the slice is empty.
func Percentile(sorted []float64, n float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if n <= 0 {
		return sorted[0]
	}
	if n >= 100 {
		return sorted[len(sorted)-1]
	}
	index := (n / 100) * float64(len(sorted)-1)
	lo := int(index)
	hi := lo + 1
	if hi >= len(sorted) {
		return sorted[lo]
	}
	frac := index - float64(lo)
	return sorted[lo] + frac*(sorted[hi]-sorted[lo])
}

// SortedDurationsMs returns a sorted slice of durations in milliseconds
// derived from the Result slice passed in.
func SortedDurationsMs(results []Result) []float64 {
	vals := make([]float64, 0, len(results))
	for _, r := range results {
		vals = append(vals, msFloat(r.Duration))
	}
	sort.Float64s(vals)
	return vals
}
