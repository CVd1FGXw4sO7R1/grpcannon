package report

import "time"

// Report holds aggregated statistics for a load test run.
type Report struct {
	Total   int
	Success int
	Failure int
	Avg     time.Duration
	Min     time.Duration
	Max     time.Duration
	P50     time.Duration
	P95     time.Duration
	P99     time.Duration
}

// Summarize builds a Report from a slice of Results.
func Summarize(results []Result) *Report {
	r := &Report{Total: len(results)}
	if len(results) == 0 {
		return r
	}

	var durations []time.Duration
	var total time.Duration
	r.Min = results[0].Duration
	r.Max = results[0].Duration

	for _, res := range results {
		if res.IsSuccess() {
			r.Success++
		} else {
			r.Failure++
		}
		durations = append(durations, res.Duration)
		total += res.Duration
		if res.Duration < r.Min {
			r.Min = res.Duration
		}
		if res.Duration > r.Max {
			r.Max = res.Duration
		}
	}

	r.Avg = total / time.Duration(len(results))
	sorted := SortedDurationsMs(durations)
	r.P50 = time.Duration(Percentile(sorted, 50)) * time.Millisecond
	r.P95 = time.Duration(Percentile(sorted, 95)) * time.Millisecond
	r.P99 = time.Duration(Percentile(sorted, 99)) * time.Millisecond
	return r
}
