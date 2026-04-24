package report

import (
	"fmt"
	"io"
	"sort"
)

// PercentileHeatmapRow holds latency data for a single percentile across time buckets.
type PercentileHeatmapRow struct {
	Percentile float64
	Buckets    []float64 // latency in ms per time bucket
}

// BuildPercentileHeatmap computes a 2D matrix of latency percentiles over time buckets.
// rows: percentile steps (e.g. 50, 75, 90, 95, 99)
// cols: time buckets of equal duration
func BuildPercentileHeatmap(results []Result, buckets int) []PercentileHeatmapRow {
	if len(results) == 0 || buckets <= 0 {
		return nil
	}

	percentiles := []float64{50, 75, 90, 95, 99}

	// partition results into time buckets
	sort.Slice(results, func(i, j int) bool {
		return results[i].StartedAt.Before(results[j].StartedAt)
	})

	start := results[0].StartedAt
	end := results[len(results)-1].StartedAt
	totalNs := end.Sub(start).Nanoseconds()
	if totalNs <= 0 {
		totalNs = 1
	}

	bucketSlices := make([][]Result, buckets)
	for i := range bucketSlices {
		bucketSlices[i] = []Result{}
	}

	for _, r := range results {
		offset := r.StartedAt.Sub(start).Nanoseconds()
		idx := int(float64(offset) / float64(totalNs) * float64(buckets))
		if idx >= buckets {
			idx = buckets - 1
		}
		bucketSlices[idx] = append(bucketSlices[idx], r)
	}

	rows := make([]PercentileHeatmapRow, len(percentiles))
	for pi, pct := range percentiles {
		row := PercentileHeatmapRow{
			Percentile: pct,
			Buckets:    make([]float64, buckets),
		}
		for bi, bucket := range bucketSlices {
			row.Buckets[bi] = Percentile(bucket, pct)
		}
		rows[pi] = row
	}
	return rows
}

// WritePercentileHeatmap writes an ASCII representation of the percentile heatmap.
func WritePercentileHeatmap(w io.Writer, results []Result, buckets int) {
	rows := BuildPercentileHeatmap(results, buckets)
	if rows == nil {
		fmt.Fprintln(w, "no data")
		return
	}

	fmt.Fprintf(w, "%-6s", "P\\T")
	for i := 0; i < buckets; i++ {
		fmt.Fprintf(w, " %6d", i+1)
	}
	fmt.Fprintln(w)

	for _, row := range rows {
		fmt.Fprintf(w, "P%-5.0f", row.Percentile)
		for _, v := range row.Buckets {
			if v == 0 {
				fmt.Fprintf(w, " %6s", "-")
			} else {
				fmt.Fprintf(w, " %6.1f", v)
			}
		}
		fmt.Fprintln(w)
	}
}
