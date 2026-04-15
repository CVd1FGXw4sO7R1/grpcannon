package report

import (
	"fmt"
	"io"
	"time"
)

const bucketCount = 8

// PrintHistogram writes an ASCII latency histogram to w.
func (s *Summary) PrintHistogram(w io.Writer) {
	if len(s.Durations) == 0 {
		fmt.Fprintln(w, "no successful requests to histogram")
		return
	}

	min := s.Durations[0]
	max := s.Durations[len(s.Durations)-1]
	span := max - min
	if span == 0 {
		span = 1
	}

	buckets := make([]int, bucketCount)
	for _, d := range s.Durations {
		idx := int(float64(d-min) / float64(span) * float64(bucketCount-1))
		if idx >= bucketCount {
			idx = bucketCount - 1
		}
		buckets[idx]++
	}

	maxCount := 0
	for _, c := range buckets {
		if c > maxCount {
			maxCount = c
		}
	}

	fmt.Fprintln(w, "\nLatency histogram:")
	barWidth := 40
	for i, c := 0, 0; i < bucketCount; i++ {
		c = buckets[i]
		lo := min + time.Duration(float64(span)*float64(i)/float64(bucketCount))
		hi := min + time.Duration(float64(span)*float64(i+1)/float64(bucketCount))
		filled := 0
		if maxCount > 0 {
			filled = c * barWidth / maxCount
		}
		bar := ""
		for j := 0; j < filled; j++ {
			bar += "█"
		}
		fmt.Fprintf(w, "  [%8s – %8s] %s (%d)\n",
			lo.Round(time.Microsecond),
			hi.Round(time.Microsecond),
			bar, c)
	}
}
