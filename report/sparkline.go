package report

import (
	"fmt"
	"io"
	"math"
	"sort"
	"time"
)

var sparkChars = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// WriteSparkline writes a latency sparkline histogram to w.
func WriteSparkline(w io.Writer, r *Report) error {
	if r == nil || len(r.Results) == 0 {
		_, err := fmt.Fprintln(w, "No results.")
		return err
	}

	durations := make([]time.Duration, 0, len(r.Results))
	for _, res := range r.Results {
		if res.Err == nil {
			durations = append(durations, res.Duration)
		}
	}
	if len(durations) == 0 {
		_, err := fmt.Fprintln(w, "No successful results.")
		return err
	}

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })

	const buckets = 16
	min := durations[0].Seconds() * 1000
	max := durations[len(durations)-1].Seconds() * 1000
	if max == min {
		max = min + 1
	}

	counts := make([]int, buckets)
	for _, d := range durations {
		ms := d.Seconds() * 1000
		idx := int(math.Floor((ms - min) / (max - min) * buckets))
		if idx >= buckets {
			idx = buckets - 1
		}
		counts[idx]++
	}

	maxCount := 0
	for _, c := range counts {
		if c > maxCount {
			maxCount = c
		}
	}

	line := make([]rune, buckets)
	for i, c := range counts {
		if maxCount == 0 {
			line[i] = sparkChars[0]
		} else {
			idx := int(math.Round(float64(c) / float64(maxCount) * float64(len(sparkChars)-1)))
			line[i] = sparkChars[idx]
		}
	}

	_, err := fmt.Fprintf(w, "Latency distribution (%.2fms – %.2fms):\n%s\n", min, max, string(line))
	return err
}
